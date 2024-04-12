package main

import (
	"log"
	"strconv"

	"github.com/albrow/zoom"
	"github.com/gin-gonic/gin"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
)

func main() {
	pool := zoom.NewPool("redis:6379")
	defer pool.Close()

	libs, err := pool.NewCollectionWithOptions(&models.Lib{}, zoom.CollectionOptions{
		Index: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	api := gin.Default()
	svr := &server{
		libs: libs,
	}
	svr.setupRoutes(api)

	if err := api.Run("0.0.0.0:80"); err != nil {
		log.Fatal(err)
	}

}

type server struct {
	libs *zoom.Collection
}

func (svr *server) setupRoutes(r gin.IRouter) {
	r.GET("/topics")
	r.GET("/topic/:name")
	r.GET("/categories")
	r.GET("/category/:name")
	r.GET("/libs", paginate, svr.getLibs)
	r.GET("/lib/:owner/:name", svr.getLib)
	r.GET("/libs/popular", paginate, svr.getLibsPopular)
	r.GET("/libs/growing", paginate, svr.getLibsGrowing)
	r.GET("/libs/recent")
	r.GET("/libs/active")
}

func paginate(c *gin.Context) {
	page := 1
	if p := c.Query("page"); p != "" {
		_p, err := strconv.Atoi(p)
		if err == nil {
			page = _p
		}
	}

	perp := 20
	if p := c.Query("per_page"); p != "" {
		_p, err := strconv.Atoi(p)
		if err == nil {
			perp = _p
		}
	}

	c.Set("limit", perp)
	c.Set("offset", (page-1)*perp)
}

func (s *server) getLib(c *gin.Context) {
	owner := c.Param("owner")
	name := c.Param("name")

	id := owner + "/" + name
	lib := models.Lib{}
	if err := s.libs.Find(id, &lib); err != nil {
		log.Printf("Error finding lib: %s", err)
		c.AbortWithError(404, err)
		return
	}

	c.JSON(200, lib)
}

func (s *server) getLibs(c *gin.Context) {
	libs := []models.Lib{}

	q := s.libs.NewQuery().
		Order("-UpdatedAt").
		Limit(c.GetUint("limit")).
		Offset(c.GetUint("offset"))

	if err := q.Run(&libs); err != nil {
		log.Printf("Error getting libs: %s", err)
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, libs)
}

func (s *server) getLibsPopular(c *gin.Context) {
	libs := []models.Lib{}

	q := s.libs.NewQuery().
		Order("-Stargazers").
		Limit(c.GetUint("limit")).
		Offset(c.GetUint("offset"))

	if err := q.Run(&libs); err != nil {
		log.Printf("Error getting popular libs: %s", err)
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, libs)
}

func (s *server) getLibsGrowing(c *gin.Context) {
	libs := []models.Lib{}

	q := s.libs.NewQuery().
		Order("-StargazersChange").
		Limit(c.GetUint("limit")).
		Offset(c.GetUint("ofsset"))

	if err := q.Run(&libs); err != nil {
		log.Printf("Error getting growing libs: %s", err)
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, libs)
}
