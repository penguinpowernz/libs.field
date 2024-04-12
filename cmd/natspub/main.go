package main

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subj := os.Args[1]
	payload := os.Args[2]

	if err := nc.Publish(subj, []byte(payload)); err != nil {
		log.Fatal(err)
	}
}
