package parsers

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
)

var binaryWordsRE = regexp.MustCompile(`[Ff]ree[Bb][Ss][Dd]|[Ll]inux|[Dd]arwin|\.exe|\.deb|\.rpm|amd64|arm64|386|ppc64|\.AppImage|\.dmg`)

// Releases will scan the releases endpoint for indications that this lib is a downloadable runnable application
func Releases(ctx context.Context, nc *nats.Conn, libs *zoom.Collection) error {
	sub, err := nc.QueueSubscribe("releases", "parsers", func(m *nats.Msg) {
		releaseParser{
			find:       libs.Find,
			save:       libs.Save,
			publish:    nc.Publish,
			saveFields: libs.SaveFields,
		}.parse(m.Data)
	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	return sub.Unsubscribe()
}

type release struct {
	URL         string `json:"url"`
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name string `json:"name"`
	} `json:"assets"`
}

type releaseParser struct {
	find       func(string, zoom.Model) error
	save       func(zoom.Model) error
	saveFields func([]string, zoom.Model) error
	publish    func(string, []byte) error
}

func (rp releaseParser) parse(data []byte) {
	releases := []release{}
	err := json.Unmarshal(data, &releases)
	if err != nil {
		log.Printf("[parsers.releases] Error parsing JSON: %s", err)
		return
	}

	if len(releases) == 0 {
		return
	}

	urlparts := strings.Split(releases[0].URL, "/")
	id := strings.Join(urlparts[4:6], "/")

	lib := &models.Lib{}
	if err := rp.find(id, lib); err != nil {
		log.Printf("[parsers.releases] Error finding lib %s: %s", id, err)
		return
	}

	if lib.ReleaseTag != releases[0].TagName {
		lib.ReleaseTag = releases[0].TagName
		lib.ReleasedAt = releases[0].PublishedAt
	}

	lib.ReleasesCheckedAt = time.Now()
	if err := rp.saveFields([]string{"ReleasesCheckedAt"}, lib); err != nil {
		log.Printf("[parsers.releases] Error updating release check time for lib %s: %s", id, err)
	}

	isApp := false
	for _, rls := range releases {
		for _, asset := range rls.Assets {
			if binaryWordsRE.MatchString(asset.Name) {
				isApp = true
				break
			}
		}
	}

	if !isApp {
		return
	}

	lib.IsApp = true
	if err := rp.save(lib); err != nil {
		log.Printf("[parsers.releases] Error saving lib %s: %s", id, err)
		return
	}

	if err := rp.publish("taxonomize", []byte(id)); err != nil {
		log.Printf("[parsers.releases] Error publishing taxonomize for lib %s: %s", id, err)
	}
}
