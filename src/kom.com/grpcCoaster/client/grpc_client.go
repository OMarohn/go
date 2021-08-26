package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"kom.com/m/v2/src/kom.com/grpcCoaster"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Konnte keine Verbindung aufbauen: %v", err)
	}
	defer conn.Close()
	log.Println("Verbindung aufgebaut.")

	// Client erstellen
	client := grpcCoaster.NewCoasterServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Println("ctx aufgebaut")

	// Abrufen der aller gespeicherten Daten
	res, err := client.GetCoasters(ctx, &grpcCoaster.Empty{})
	if err != nil {
		log.Fatalf("Fehler beim Aufruf: %v", err)
	}
	log.Printf("Result: %v", res)

	// Anlegen eines Datums
	coaster := grpcCoaster.CoasterMessage{}
	coaster.Id = "id88"
	coaster.Height = 88
	coaster.Name = "Neu angelegt - grpc"
	coaster.Manufacture = "Google"

	_, err = client.CreateCoaster(ctx, &coaster)
	if err != nil {
		log.Printf("Fehler beim Aufruf: %v", err)
	}

	// Abfragen eines Datums by ID
	res2, err := client.GetCoaster(ctx, &grpcCoaster.CoasterIDMessage{Id: "id88"})
	if err != nil {
		log.Printf("Fehler beim Aufruf: %v", err)
	}
	log.Printf("Resultat: %v", res2)
}
