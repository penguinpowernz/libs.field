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

func Commits(ctx context.Context, nc *nats.Conn, libs *zoom.Collection) error {
	sub, err := nc.QueueSubscribe("commits", "parsers", func(m *nats.Msg) {
		log.Printf("[parsers.commits] got message")
		commitsParser{
			find:       libs.Find,
			save:       libs.Save,
			saveFields: libs.SaveFields,
		}.parse(m.Data)
	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	return sub.Unsubscribe()
}

type commit struct {
	URL    string `json:"url"` // https://api.github.com/repos/gin-gonic/gin/commits/0397e5e0c0f8f8176c29f7edd8f1bff8e45df780
	Commit struct {
		Author struct {
			Date string `json:"date"` // 2024-04-07T02:18:23Z
		} `json:"author"`
	} `json:"commit"`
}

type commitsParser struct {
	save       func(zoom.Model) error
	find       func(string, zoom.Model) error
	saveFields func([]string, zoom.Model) error
}

func (cp commitsParser) parse(data []byte) {
	commits := []commit{}
	err := json.Unmarshal(data, &commits)
	if err != nil {
		log.Printf("[parsers.commits] Error parsing JSON: %s", err)
		return
	}

	if len(commits) == 0 {
		return
	}

	urlparts := strings.Split(commits[0].URL, "/")
	id := strings.Join(urlparts[4:6], "/")

	lib := &models.Lib{}
	if err := cp.find(id, lib); err != nil {
		log.Printf("[parsers.commits] Error finding lib %s: %s", id, err)
		return
	}

	lib.CommitsCheckedTime = time.Now()
	if err := cp.saveFields([]string{"CommitsCheckedTime"}, lib); err != nil {
		log.Printf("[parsers.commits] Error updating commit check time for lib %s: %s", id, err)
	}

	days := map[string]int{}
	for _, commit := range commits {
		day := strings.Split(commit.Commit.Author.Date, "T")[0]
		days[day]++
	}

	avg := 0
	for _, count := range days {
		avg += count
	}

	avg /= len(days)

	lib.PushesPerday = avg

	if err := cp.save(lib); err != nil {
		log.Printf("[parsers.commits] Error saving lib %s: %s", id, err)
	}
}
