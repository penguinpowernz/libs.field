package parsers

import (
	"context"
	"log"

	"github.com/Jeffail/gabs"
	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
)

func Search(ctx context.Context, nc *nats.Conn, _ *zoom.Collection) error {
	sub, err := nc.QueueSubscribe("search", "parsers", func(m *nats.Msg) {
		err := parseSearch(m.Data, func(data []byte) {
			// call this for each repo item found in the search results
			if err := nc.Publish("repos", data); err != nil {
				log.Printf("[parsers.search] Error publishing repos: %s", err)
			}
		}, func(string) {})

		if err != nil {
			log.Printf("[parsers.search] Error parsing search: %s", err)
		}
	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	return sub.Unsubscribe()
}

func parseSearch(data []byte, repocb func([]byte), urlcb func(string)) error {
	c, err := gabs.ParseJSON(data)
	if err != nil {
		return err
	}

	items, err := c.Path("items").Children()
	if err != nil {
		return err
	}

	for _, item := range items {
		repocb(item.Bytes())
	}

	return nil
}
