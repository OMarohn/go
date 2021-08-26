package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"kom.com/m/v2/src/kom.com/coaster/coaster"
	"kom.com/m/v2/src/kom.com/grpcCoaster"
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
	log.Println("REDIS connected")

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to Listen on Port 9000; %v", err)
	}

	port := coaster.NewCoasterGrpcServerPort(
		coaster.NewCoasterService(
			coaster.NewRedisRepo(conn.RedisClient)))

	grpcServer := grpc.NewServer()

	grpcCoaster.RegisterCoasterServiceServer(grpcServer, &port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to create grpc-Server on Port 9000; %v", err)
	}
}
