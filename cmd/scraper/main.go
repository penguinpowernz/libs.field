package main

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/scraper"
)

func main() {
	nc, err := nats.Connect("nats:4222")
	if err != nil {
		log.Printf("Error connecting to nats: %s", err)
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing nats connection")
		nc.Close()
	}()

	s := scraper.New(nc)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		log.Println("Stopping scraper")
		cancel()
	}()

	log.Println("Starting scraper")
	s.Run(ctx)
}
