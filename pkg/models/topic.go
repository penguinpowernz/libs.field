package models

import "github.com/albrow/zoom"

type Topic struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (t *Topic) ModelID() string {
	return t.Name
}

func (t *Topic) SetModelID(id string) {
	t.Name = id
}

type LibTopic struct {
	Lib   string `json:"lib" zoom:"index"`
	Topic string `json:"topic" zoom:"index"`
	zoom.RandomID
}
