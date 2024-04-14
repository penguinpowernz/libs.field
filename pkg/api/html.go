package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
)

func (s *Server) index(c *gin.Context) {
	sort := c.Query("sort")

	view := gin.H{}

	if count, err := s.libs.NewQuery().Count(); err == nil {
		view["count"] = count
		view["pagecount"] = count/int(c.GetUint("limit")) + 1
	}

	var libs []*models.Lib
	q := s.libs.NewQuery().
		Limit(c.GetUint("limit")).
		Offset(c.GetUint("offset"))

	switch sort {
	case "popular":
		view["sort"] = "popular"
		q.Order("-Stargazers")

	case "pushed":
		q.Order("-PushedAt")
		view["sort"] = "pushed"
	case "active":
		q.Order("-PushesPerday")
		view["sort"] = "active"
	case "growing":
		q.Order("-StargazersChange")
		view["sort"] = "growing"
	case "released":
		q.Order("-ReleasedAt")
		view["sort"] = "released"
	default:

	}

	if err := q.Run(&libs); err != nil {
		log.Printf("Error getting popular libs: %s", err)
		c.AbortWithError(500, err)
		return
	}
	view["libs"] = libs

	var cats []*models.Category
	if err := s.cats.FindAll(&cats); err != nil {
		log.Printf("Error getting categories: %s", err)
		c.AbortWithError(500, err)
		return
	}
	view["cats"] = cats

	c.HTML(200, "index.html", view)
}
