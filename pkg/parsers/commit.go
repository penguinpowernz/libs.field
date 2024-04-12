package parsers

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
)

// Commit is for updated the tagged at value for a commit
func Commit(ctx context.Context, nc *nats.Conn, libs *zoom.Collection) error {
	sub, err := nc.QueueSubscribe("commit", "parsers", func(m *nats.Msg) {
		commitParser{
			find:       libs.Find,
			saveFields: libs.SaveFields,
		}.parse(m.Data)
	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	return sub.Unsubscribe()
}

type commitParser struct {
	find       func(string, zoom.Model) error
	saveFields func([]string, zoom.Model) error
}

func (cp commitParser) parse(data []byte) {
	c := commit{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		log.Printf("[parsers.commit] Error parsing JSON: %s", err)
		return
	}

	urlparts := strings.Split(c.URL, "/")
	id := strings.Join(urlparts[4:6], "/")

	log.Printf("[parsers.commit] Processing commit for lib %s", id)

	lib := &models.Lib{}
	if err := cp.find(id, lib); err != nil {
		log.Printf("[parsers.commit] Error finding lib %s: %s", id, err)
		return
	}

	lib.TaggedAt = c.Commit.Author.Date
	if err := cp.saveFields([]string{"TaggedAt"}, lib); err != nil {
		log.Printf("[parsers.commit] Error updating lib %s: %s", id, err)
		return
	}
}
