// https://archium.org/index.php/Der_Golang-Spicker#Ein_Go-Spickzettel_.28.22Cheat_Sheet.22.29
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		token, found := r.Header["Authorization"]
		if found {
			log.Println(token)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

/// JWT Stuff

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

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://dev-vdt9zz3q.us.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func main() {
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

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			/**
			// Verify 'aud' claim
			aud := "YOUR_API_IDENTIFIER"
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}

			// Verify 'iss' claim
			iss := "https://YOUR_DOMAIN/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}
			**/

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

			return result, err
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodRS256,
		Debug:         true,
	})

	r := mux.NewRouter()
	r.Use(jwtMiddleware.Handler)

	r.HandleFunc("/coasters", port_REST_mem.HandleList).Methods(http.MethodGet)
	r.HandleFunc("/coasters/{id}", port_REST_mem.HandleGetOne).Methods(http.MethodGet)
	r.HandleFunc("/coasters", port_REST_mem.HandleCreate).Methods(http.MethodPost)

	sr := r.PathPrefix("/redis").Subrouter()
	// sr.Use(loggingMiddleware)
	sr.HandleFunc("/coasters", port_REST_redis.HandleList).Methods(http.MethodGet)
	sr.HandleFunc("/coasters/{id}", port_REST_redis.HandleGetOne).Methods(http.MethodGet)
	sr.HandleFunc("/coasters/{id}", port_REST_redis.HandleDelete).Methods(http.MethodDelete)
	sr.HandleFunc("/coasters", port_REST_redis.HandleCreate).Methods(http.MethodPost)

	sr.HandleFunc("/extern", func(w http.ResponseWriter, r *http.Request) {
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

	}).Methods(http.MethodGet)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
