package main

import (
	"context"
	"log"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
	"github.com/penguinpowernz/libs.fieid/pkg/taxon"
)

func main() {
	nc, err := nats.Connect("nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	pool := zoom.NewPool("redis:6379")
	defer pool.Close()

	libs, err := pool.NewCollectionWithOptions(&models.Lib{}, zoom.CollectionOptions{
		Index: true,
	})

	tmzr := taxon.New(nc, libs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tmzr.Run(ctx)
}
