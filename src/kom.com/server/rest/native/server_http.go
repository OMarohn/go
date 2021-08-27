// https://archium.org/index.php/Der_Golang-Spicker#Ein_Go-Spickzettel_.28.22Cheat_Sheet.22.29
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
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

func main() {
	conn, err := RedisClient()
	defer conn.closeRedis()

	if err != nil {
		panic(err)
	}

	// Daten aus dem REDIS
	port := coaster.NewCoasterRestPort(
		coaster.NewCoasterService(
			coaster.NewRedisRepo(conn.RedisClient)))

	// Daten aus dem Speicher
	port_mem := coaster.NewCoasterRestPort(
		coaster.NewCoasterService(
			coaster.NewCoasterMemmoryRepo()))

	// Handler
	http.HandleFunc("/coasters", port_mem.Handle)
	http.HandleFunc("/redis/coasters", port.Handle)
	http.HandleFunc("/redis/coasters/", port.HandleGet)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
