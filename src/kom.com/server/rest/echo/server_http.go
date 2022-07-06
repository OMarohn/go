package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	_ "github.com/lib/pq"
	"kom.com/m/v2/src/kom.com/coaster/coaster"
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

func main() {
	log.Println("Echo Server")

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

	// Daten aus dem Speicher
	port_REST_mem := coaster.NewCoasterRestPort3(
		coaster.NewCoasterService(
			coaster.NewCoasterMemmoryRepo()))

	// port_REST_db := coaster.NewCoasterRestPort2(
	// 	coaster.NewCoasterService(
	// 		coaster.NewPostgresRepo(conn.CoasterDB)))

	/**
	Ab hier ECHO Server und Middleware
	*/

	/**
	Routing
	*/
	e := echo.New()

	e.GET("/coasters", port_REST_mem.HandleList)
	e.GET("/coasters/:id", port_REST_mem.HandleGetOne)
	e.POST("/coasters", port_REST_mem.HandleCreate)

	sr := e.Group("/redis")
	sr.DELETE("/coasters/:id", port_REST_redis.HandleDelete)
	sr.GET("/coasters", port_REST_redis.HandleList)

	sr.GET("/coasters/:id", port_REST_redis.HandleGetOne)
	sr.POST("/coasters", port_REST_redis.HandleCreate)

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
