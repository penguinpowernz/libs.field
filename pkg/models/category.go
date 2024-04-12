package models

import "github.com/albrow/zoom"

type Category struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (c *Category) ModelID() string {
	return c.Name
}

func (c *Category) SetModelID(id string) {
	c.Name = id
}

type LibCategory struct {
	Lib      string `json:"lib" zoom:"index"`
	Category string `json:"category" zoom:"index"`
	zoom.RandomID
}
