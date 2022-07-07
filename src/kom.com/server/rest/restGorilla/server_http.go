// https://archium.org/index.php/Der_Golang-Spicker#Ein_Go-Spickzettel_.28.22Cheat_Sheet.22.29
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/lib/pq"
	"kom.com/m/v2/src/kom.com/coaster/coaster"
	jwkst "kom.com/m/v2/src/kom.com/server/jwks"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "gocoaster_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// Connections
type Connections struct {
	RedisClient *redis.Client
	CoasterDB   *sql.DB
}

// Erzeugen der Clients fuer die DB und den REDIS
func ConnectionClient() (*Connections, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPsw := os.Getenv("REDIS_PSW")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPsw,
		DB:       0,
	})

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPsw := os.Getenv("DB_PSW")
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPsw, "coaster", dbHost, dbPort)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		panic("Konnte DB nich öffnen!")
	}

	log.Println(dbinfo)
	log.Println(db)

	return &Connections{
		RedisClient: rdb,
		CoasterDB:   db,
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

// Prometheus Metriken
// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

func main() {
	log.Println("NEU4")
	jwkManager := jwkst.JWKManager{Url: "https://dev-vdt9zz3q.us.auth0.com/.well-known/jwks.json", Filename: "/app/jwks.json"}
	jwkManager.InitCertStore()

	conn, err := ConnectionClient()
	// defer auch noch DB schliessen
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

	port_REST_db := coaster.NewCoasterRestPort2(
		coaster.NewCoasterService(
			coaster.NewPostgresRepo(conn.CoasterDB)))

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
			pk, err := jwkManager.GetPemCert(token)
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
			pk, err := jwkManager.GetPemCert(token)
			if err != nil {
				return nil, errors.New("token invalid")
			}

			return pk, err
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	r := mux.NewRouter()
	r.Use(prometheusMiddleware)

	r.Path("/metrics").Handler(promhttp.Handler())

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

	srdb := r.PathPrefix("/db").Subrouter()
	srdb.Use(loggingMiddleware)
	srdb.HandleFunc("/coasters", port_REST_db.HandleList).Methods(http.MethodGet)
	srdb.HandleFunc("/coasters", port_REST_db.HandleCreate).Methods(http.MethodPost)
	srdb.HandleFunc("/coasters/{id}", port_REST_db.HandleGetOne).Methods(http.MethodGet)
	srdb.HandleFunc("/coasters/{id}", port_REST_db.HandleDelete).Methods(http.MethodDelete)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
