package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	jwtft "github.com/form3tech-oss/jwt-go"
	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-redis/redis/v8"
	promMW "github.com/labstack/echo-contrib/prometheus"

	_ "github.com/lib/pq"
	"kom.com/m/v2/src/kom.com/coaster/coaster"
	"kom.com/m/v2/src/kom.com/graph/generated"
	jwkstools "kom.com/m/v2/src/kom.com/server/jwks"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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

func authmw() echo.MiddlewareFunc {
	jwkManager := jwkstools.JWKManager{Url: "https://dev-vdt9zz3q.us.auth0.com/.well-known/jwks.json", Filename: "/app/jwks.json"}
	jwkManager.InitCertStore()

	// initialize JWT middleware instance
	return middleware.JWTWithConfig(middleware.JWTConfig{
		// skip public endpoints
		Skipper: func(context echo.Context) bool {
			return strings.Contains(context.Path(), "/metrics")
		},
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			// casting the hard way @todo castin the easy way lernen ;-)
			t := jwtft.Token{Raw: token.Raw, Method: token.Method, Header: token.Header, Claims: token.Claims, Signature: token.Signature, Valid: token.Valid}

			pk, err := jwkManager.GetPemCert(&t)
			if err != nil {
				return nil, errors.New("token invalid")
			}
			return pk, err
		},
	})
}

func healthHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}

func readyHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}

func CreateEchoServer() *echo.Echo {

	conn, err := ConnectionClient()
	// defer auch noch DB schliessen
	defer conn.closeRedis()

	if err != nil {
		panic(err)
	}

	// Daten aus dem REDIS
	port_REST_redis := coaster.NewCoasterRestPort3(
		coaster.NewCoasterService(
			coaster.NewRedisRepo(conn.RedisClient)))

	// Service der auf dem Arbeitsspeicher arbeitet - soll auch für den GQL-Port Verwendet werden.
	memService := coaster.NewCoasterService(
		coaster.NewCoasterMemmoryRepo())
	// Daten aus dem Speicher
	port_REST_mem := coaster.NewCoasterRestPort3(memService)

	port_REST_db := coaster.NewCoasterRestPort3(
		coaster.NewCoasterService(
			coaster.NewPostgresRepo(conn.CoasterDB)))

	/**
	Ab hier ECHO Server und Middleware
	*/
	e := echo.New()

	prom := promMW.NewPrometheus("echo", nil)
	prom.Use(e)
	e.Use(middleware.Logger())
	jwksMW := authmw()

	// gql
	resolver := coaster.NewCoasterResolver(&memService)
	gQLSrv := echo.WrapHandler(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver})))
	gQLPlayGroundHandler := echo.WrapHandler(playground.Handler("GraphQL playground", "/query"))
	e.GET("/playground", gQLPlayGroundHandler)
	e.POST("/query", gQLSrv)

	// REST
	sm := e.Group("/mem")
	sm.GET("/coasters", port_REST_mem.HandleList)
	sm.GET("/coasters/:id", port_REST_mem.HandleGetOne)
	sm.POST("/coasters", port_REST_mem.HandleCreate)
	sm.DELETE("/coasters/:id", port_REST_mem.HandleDelete)

	sr := e.Group("/redis", jwksMW)
	sr.GET("/coasters", port_REST_redis.HandleList)
	sr.GET("/coasters/:id", port_REST_redis.HandleGetOne)
	sr.POST("/coasters", port_REST_redis.HandleCreate)
	sr.DELETE("/coasters/:id", port_REST_redis.HandleDelete)

	srd := e.Group("/db")
	srd.GET("/coasters", port_REST_db.HandleList)
	srd.GET("/coasters/:id", port_REST_db.HandleGetOne)
	srd.POST("/coasters", port_REST_db.HandleCreate)
	srd.DELETE("/coasters/:id", port_REST_db.HandleDelete)

	e.GET("/healthz", healthHandler)
	e.GET("/readyz", readyHandler)

	return e

}

func main() {

	e := CreateEchoServer()

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
