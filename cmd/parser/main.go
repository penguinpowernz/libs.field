package main

import (
	"context"
	"log"
	"time"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
	"github.com/penguinpowernz/libs.fieid/pkg/parsers"
)

func main() {
	nc, err := nats.Connect("nats:4222")
	if err != nil {
		log.Fatalf("Error connecting to nats: %s", err)
	}
	defer nc.Close()

	pool := zoom.NewPool("redis:6379")
	defer pool.Close()

	libs, err := pool.NewCollectionWithOptions(&models.Lib{}, zoom.CollectionOptions{
		Index: true,
	})

	if err != nil {
		log.Fatalf("Error creating libs collection: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for name, fn := range parsers.Parsers {
		go func(name string, fn parsers.ParserFunc) {
			start := time.Now()
			log.Printf("Starting to parse %s", name)
			err := fn(ctx, nc, libs)
			if err != nil {
				log.Printf("Error parsing %s: %s (took %s)", name, err, time.Since(start))
			} else {
				log.Printf("Finished parsing %s (took %s)", name, time.Since(start))
			}
		}(name, fn)
	}

	<-ctx.Done()
}
