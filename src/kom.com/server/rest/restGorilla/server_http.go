// https://archium.org/index.php/Der_Golang-Spicker#Ein_Go-Spickzettel_.28.22Cheat_Sheet.22.29
package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/json"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"kom.com/m/v2/src/kom.com/coaster/coaster"
)

// Connections
type Connections struct {
	RedisClient *redis.Client
}

func RedisClient() (*Connections, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPsw := os.Getenv("REDIS_PSW")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPsw,
		DB:       0,
	})

	return &Connections{
		RedisClient: rdb,
	}, nil
}

// Schliessen der Connections -> Connections
func (d *Connections) closeRedis() error {
	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing Redis Client: %w", err)
	}
	return nil
}

// Einfache Middelware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		next.ServeHTTP(w, r)
	})
}

/// JWT Stuff

// Struktur des JWKS
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// Remote lesen des JWKS
func getJWKS(url string, certStore map[string]*rsa.PublicKey) (int, error) {
	resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return 0, err
	}

	return mapJwks2Store(jwks, certStore)
}

// JWKS Ergbniss in eine Map id|rsaPublic Key wandeln
func mapJwks2Store(items Jwks, certStore map[string]*rsa.PublicKey) (int, error) {
	for k, _ := range items.Keys {
		cert := "-----BEGIN CERTIFICATE-----\n" + items.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			log.Println("Public Key nicht gefunden")
		}
		certStore[items.Keys[k].Kid] = pubKey
	}

	return len(certStore), nil
}

// Public Key ermitteln
func getPemCert(token *jwt.Token, certStore map[string]*rsa.PublicKey) (*rsa.PublicKey, error) {

	kid := fmt.Sprintf("%v", token.Header["kid"])
	pk, found := certStore[kid]
	if !found {
		err := errors.New("konnte key für das zertifikat nicht finden")
		return nil, err
	}

	return pk, nil
}

func checkScope(claims jwt.MapClaims, scope string) bool {
	ret := false

	scopes := fmt.Sprintf("%v", claims["scope"])
	log.Println(scopes)
	result := strings.Split(scopes, " ")
	for i := range result {
		if result[i] == scope {
			ret = true
			break
		}
	}

	return ret
}

func main() {
	certStore := map[string]*rsa.PublicKey{}
	// JWT Certifikate preloaden -- hier hab ich noch nen Problem im ISTIO-EGress!
	cnt, err := getJWKS("https://dev-vdt9zz3q.us.auth0.com/.well-known/jwks.json", certStore)
	if err != nil {
		log.Println("Keine Zerifikate im Store gefunden!", err)
		// Datei mit JWKS lesen
		jsonFile, err := os.Open("./jwks.json")

		if err != nil {
			panic("Auch keine JWKS-Datei gefunden!")
		} else {
			defer jsonFile.Close()
			log.Println("File OK")
			var jwks = Jwks{}
			err = json.NewDecoder(jsonFile).Decode(&jwks)
			if err != nil {
				panic(err)
			} else {
				cnt, _ = mapJwks2Store(jwks, certStore)
			}
		}
	}
	log.Println("Anzahl Zerifikate: ", cnt)

	conn, err := RedisClient()
	defer conn.closeRedis()

	if err != nil {
		panic(err)
	}

	// Daten aus dem REDIS
	port_REST_redis := coaster.NewCoasterRestPort2(
		coaster.NewCoasterService(
			coaster.NewRedisRepo(conn.RedisClient)))

	// Daten aus dem Speicher
	port_REST_mem := coaster.NewCoasterRestPort2(
		coaster.NewCoasterService(
			coaster.NewCoasterMemmoryRepo()))

	// JWT checken
	jwtMiddlewareRO := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			mapClaims := token.Claims.(jwt.MapClaims)
			// Uhr läuft auseinander, deshalb kein Check aus den Ausstellungszeitpunkt
			delete(mapClaims, "iat")

			// Token sollte noch nicht abgelaufen sein
			checkValid := mapClaims.VerifyExpiresAt(time.Now().Unix(), true)
			if !checkValid {
				return nil, errors.New("token outdated")
			}

			if !checkScope(mapClaims, "read:sample") {
				return nil, errors.New("scope not supported")
			}

			// das mit dem CertStore ist noch ein wenig unschön!
			pk, err := getPemCert(token, certStore)
			if err != nil {
				return nil, errors.New("token invalid")
			}

			return pk, err
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	jwtMiddlewareRW := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			mapClaims := token.Claims.(jwt.MapClaims)
			// Uhr läuft auseinander, deshalb kein Check aus den Ausstellungszeitpunkt
			delete(mapClaims, "iat")

			// Token sollte noch nicht abgelaufen sein
			checkValid := mapClaims.VerifyExpiresAt(time.Now().Unix(), true)
			if !checkValid {
				return nil, errors.New("token outdated")
			}

			if !checkScope(mapClaims, "write:sample") {
				return nil, errors.New("scope not supported")
			}

			// das mit dem CertStore ist noch ein wenig unschön!
			pk, err := getPemCert(token, certStore)
			if err != nil {
				return nil, errors.New("token invalid")
			}

			return pk, err
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	r := mux.NewRouter()

	r.HandleFunc("/coasters", port_REST_mem.HandleList).Methods(http.MethodGet)
	r.HandleFunc("/coasters/{id}", port_REST_mem.HandleGetOne).Methods(http.MethodGet)
	r.HandleFunc("/coasters", port_REST_mem.HandleCreate).Methods(http.MethodPost)

	sr := r.PathPrefix("/redis").Subrouter()
	srg := sr.Methods(http.MethodGet).Subrouter()
	srg.Use(jwtMiddlewareRO.Handler)

	srp := sr.Methods(http.MethodPost).Subrouter()
	srp.Use(jwtMiddlewareRW.Handler)

	sr.HandleFunc("/coasters/{id}", port_REST_redis.HandleDelete).Methods(http.MethodDelete)

	srg.HandleFunc("/coasters", port_REST_redis.HandleList)
	srg.HandleFunc("/coasters/{id}", port_REST_redis.HandleGetOne)
	srp.HandleFunc("/coasters", port_REST_redis.HandleCreate)

	srg.HandleFunc("/extern", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")

		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatalln(err)
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		w.Write(body)

	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
