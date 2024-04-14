package main

import (
	"log"
	"time"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
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

	libs, err := pool.NewCollectionWithOptions(&models.Lib{}, zoom.CollectionOptions{
		FallbackMarshalerUnmarshaler: util.FallbackMarshaler{},
		Index:                        true,
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		allLibs := []*models.Lib{}
		if err := libs.FindAll(&allLibs); err != nil {
			log.Println(err)
			time.Sleep(time.Second * 10)
			continue
		}

		for _, lib := range allLibs {
			lag := time.Since(lib.UpdatedAt).Seconds()
			if lag > float64(86400) {
				nc.Publish("urls", []byte(lib.APIURL))
			}
		}

		time.Sleep(time.Minute * 30)
	}

}
