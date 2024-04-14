package taxon

import (
	"context"
	"log"
	"strings"

	"github.com/albrow/zoom"
	"github.com/nats-io/nats.go"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
)

func New(nc *nats.Conn, libs, cats, libcats, topics, libtopics *zoom.Collection) *Taxonomizer {
	t := &Taxonomizer{
		nc:            nc,
		libs:          libs,
		categories:    cats,
		libCategories: libcats,
		topics:        topics,
		libTopics:     libtopics,
	}

	return t
}

type Taxonomizer struct {
	nc            *nats.Conn
	libs          *zoom.Collection
	topics        *zoom.Collection
	categories    *zoom.Collection
	libCategories *zoom.Collection
	libTopics     *zoom.Collection
}

func (t *Taxonomizer) SetupDefaults(ctx context.Context) {
	for _, name := range defaultCategories {
		exists, err := t.categories.Exists(name)
		if err != nil {
			log.Printf("Error checking if category exists: %s", err)
			continue
		}

		if !exists {
			if err := t.categories.Save(&models.Category{Name: name}); err != nil {
				log.Printf("Error inserting category: %s", err)
				continue
			}
		}
	}
}

func (t *Taxonomizer) UpdateCategoryCounts() {
	for _, name := range defaultCategories {
		count, err := t.libCategories.NewQuery().Filter("Category = ", name).Count()
		if err != nil {
			log.Printf("Error counting lib_categories: %s", err)
			continue
		}

		c := models.Category{
			Name:  name,
			Count: count,
		}

		if err := t.categories.Save(&c); err != nil {
			log.Printf("Error updating category: %s", err)
			continue
		}
	}
}

func (t *Taxonomizer) UpdateTopicCounts() {
	topics := []models.Topic{}

	if err := t.topics.FindAll(&topics); err != nil {
		log.Printf("Error getting topics: %s", err)
		return
	}

	for _, topic := range topics {
		count, err := t.libTopics.NewQuery().Filter("Topic = ", topic.Name).Count()
		if err != nil {
			log.Printf("Error counting lib_topics: %s", err)
			continue
		}

		topic.Count = count
		if err := t.topics.Save(&topic); err != nil {
			log.Printf("Error updating topic: %s", err)
			continue
		}
	}
}

func (t *Taxonomizer) Run(ctx context.Context) error {
	sub, err := t.nc.QueueSubscribe("taxonomizer", "taxon", func(m *nats.Msg) {
		bits := strings.Split(string(m.Data), " ")
		id := bits[0]
		topics := []string{}

		if len(bits) > 1 {
			topics = strings.Split(bits[1], ",")
		}

		lib := &models.Lib{}
		if err := t.libs.Find(id, lib); err != nil {
			log.Printf("Error finding lib: %s", err)
			return
		}

		for _, topic := range topics {
			count, err := t.libTopics.NewQuery().Filter("Lib = ", id).Filter("Topic = ", topic).Count()
			if err != nil {
				log.Printf("Error counting lib_topics: %s", err)
				continue
			}

			if count == 0 {
				lt := &models.LibTopic{
					Lib:   id,
					Topic: topic,
				}
				if err := t.libTopics.Save(lt); err != nil {
					log.Printf("Error inserting lib_topic: %s", err)
				}
			}
		}

		if lib.IsApp {
			count, err := t.libCategories.NewQuery().Filter("Lib = ", id).Filter("Category = ", "Applications").Count()
			if err != nil {
				log.Printf("Error counting lib_categories: %s", err)
				return
			}

			if count == 0 {
				lt := &models.LibCategory{
					Lib:      id,
					Category: "Applications",
				}
				if err := t.libCategories.Save(lt); err != nil {
					log.Printf("Error inserting lib_category: %s", err)
				}
			}
		}

	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	return sub.Unsubscribe()
}
