package parsers

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
)

func Repos(ctx context.Context, nc *nats.Conn, libs *zoom.Collection) error {
	sub, err := nc.QueueSubscribe("repos", "parsers", func(m *nats.Msg) {
		log.Printf("[parsers.repos] got message")
		rp := &repoParser{
			find:    libs.Find,
			save:    libs.Save,
			publish: nc.Publish,
			exists:  libs.Exists,
		}

		rp.parse(m.Data)
	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	return sub.Unsubscribe()
}

type repoParser struct {
	find    func(string, zoom.Model) error
	save    func(zoom.Model) error
	publish func(string, []byte) error
	exists  func(string) (bool, error)

	repo models.GitHubRepo
	lib  *models.Lib
}

func (rp *repoParser) parse(data []byte) {
	err := json.Unmarshal(data, &rp.repo)
	if err != nil {
		log.Printf("[parsers.repos] Error parsing JSON: %s", err)
		return
	}

	// ignore forked repos
	if rp.repo.Fork {
		return
	}

	exists, err := rp.exists(rp.repo.FullName)
	if err != nil {
		log.Printf("[parsers.repos] Error checking if lib exists: %s (repo: %s)", err, rp.repo.FullName)
		return
	}

	if !exists {
		rp.create()
		return
	}

	rp.update()
}

func (rp *repoParser) delegate() {
	// update tags every day
	if int(time.Since(rp.lib.TagsCheckedTime).Hours()) > 24 {
		rp.publish("urls", []byte(rp.repo.TagsURL))
	}

	// update contributors every 2-5 days
	if int(time.Since(rp.lib.ContributorsCheckedTime).Hours()) > 24*(rand.Intn(5)+2) {
		rp.publish("urls", []byte(rp.repo.ContributorsURL))
	}

	// check if its an app every 2-5 days
	if int(time.Since(rp.lib.ReleasesCheckedTime).Hours()) > 24*(rand.Intn(5)+2) {
		rp.publish("urls", []byte(strings.Split(rp.repo.ReleasesURL, "{")[0]))
	}

	// check the commits per day every day
	if int(time.Since(rp.lib.CommitsCheckedTime).Hours()) > 24 {
		rp.publish("urls", []byte(strings.Split(rp.repo.CommitsURL, "{")[0]))
	}

	// always taxonomize
	rp.publish("taxonomize", []byte(rp.repo.FullName+" "+strings.Join(rp.repo.Topics, ",")))
}

func (rp *repoParser) create() {
	rp.lib = models.NewLibFromRepo(rp.repo)
	if err := rp.save(rp.lib); err != nil {
		log.Printf("[parsers.repos] Error saving lib: %s (repo: %s)", err, rp.repo.FullName)
		return
	}

	rp.delegate()
}

func (rp *repoParser) update() {
	rp.lib = new(models.Lib)
	if err := rp.find(rp.repo.FullName, rp.lib); err != nil {
		log.Printf("[parsers.repos] Error finding lib: %s (repo: %s)", err, rp.repo.FullName)
		return
	}

	rp.lib.UpdateFromRepo(rp.repo)

	if err := rp.save(rp.lib); err != nil {
		log.Printf("[parsers.repos] Error updating lib: %s (repo: %s)", err, rp.repo.FullName)
		return
	}

	rp.delegate()
}
