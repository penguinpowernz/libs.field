package parsers

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
)

func Tags(ctx context.Context, nc *nats.Conn, libs *zoom.Collection) error {
	sub, err := nc.QueueSubscribe("tags", "parsers", func(m *nats.Msg) {
		tagParser{
			find:       libs.Find,
			save:       libs.Save,
			saveFields: libs.SaveFields,
			publish:    nc.Publish,
		}.parse(m.Data)
	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	return sub.Unsubscribe()
}

type tag struct {
	Name   string `json:"name"`
	Commit struct {
		URL string `json:"url"`
	} `json:"commit"`
}

type tagParser struct {
	find       func(string, zoom.Model) error
	save       func(zoom.Model) error
	saveFields func([]string, zoom.Model) error
	publish    func(string, []byte) error
}

func (tp tagParser) parse(data []byte) {
	tags := []tag{}
	err := json.Unmarshal(data, &tags)
	if err != nil {
		log.Printf("[parsers.tags] Error parsing JSON: %s", err)
		return
	}

	if len(tags) == 0 {
		return
	}

	urlparts := strings.Split(tags[0].Commit.URL, "/")
	id := strings.Join(urlparts[4:6], "/")

	lib := &models.Lib{}
	if err := tp.find(id, lib); err != nil {
		log.Printf("[parsers.tags] Error finding lib %s: %s", id, err)
		return
	}

	lib.TagsCheckedAt = time.Now()
	if err := tp.saveFields([]string{"TagsCheckedAt"}, lib); err != nil {
		log.Printf("[parsers.tags] Error updating tag check time: %s", err)
	}

	if lib.CurrentTag != tags[0].Name {
		lib.CurrentTag = tags[0].Name
		lib.TaggedAt = ""
	}

	if lib.ReleaseTag == lib.CurrentTag {
		lib.TaggedAt = lib.ReleasedAt
	}

	if lib.TaggedAt == "" && lib.ReleaseTag == "" {
		tp.publish("urls", []byte(tags[0].Commit.URL))
	}

	if err := tp.save(lib); err != nil {
		log.Printf("[parsers.tags] Error updating lib: %s", err)
		return
	}
}
