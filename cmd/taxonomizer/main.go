package main

import (
	"context"
	"log"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
	"github.com/penguinpowernz/libs.fieid/pkg/taxon"
	"github.com/penguinpowernz/libs.fieid/pkg/util"
)

func main() {
	nc, err := nats.Connect("nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	pool := zoom.NewPool("redis:6379")
	defer pool.Close()

	opts := zoom.CollectionOptions{
		FallbackMarshalerUnmarshaler: util.FallbackMarshaler{},
		Index:                        true,
	}
	libs, err := pool.NewCollectionWithOptions(&models.Lib{}, opts)
	if err != nil {
		log.Fatal(err)
	}

	tmzr := taxon.New(nc, libs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tmzr.Run(ctx)
}
