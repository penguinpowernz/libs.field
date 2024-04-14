package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
)

func (s *Server) getLib(c *gin.Context) {
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

func (s *Server) getLibs(c *gin.Context) {
	libs := []*models.Lib{}

	q := s.libs.NewQuery().
		Limit(c.GetUint("limit")).
		Offset(c.GetUint("offset"))

	if err := q.Run(&libs); err != nil {
		log.Printf("Error getting libs: %s", err)
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, libs)
}

func (s *Server) getLibsPopular(c *gin.Context) {
	libs := []*models.Lib{}

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

func (s *Server) getLibsGrowing(c *gin.Context) {
	libs := []*models.Lib{}

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

func (s *Server) getCategories(c *gin.Context) {
	cats := []*models.Category{}

	if err := s.cats.FindAll(&cats); err != nil {
		log.Printf("Error getting categories: %s", err)
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, cats)
}

func (s *Server) getCategory(c *gin.Context) {
	name := c.Param("name")

	cat := models.Category{}
	if err := s.cats.Find(name, &cat); err != nil {
		log.Printf("Error finding category: %s", err)
		c.AbortWithError(404, err)
		return
	}

	libcats := []*models.LibCategory{}
	err := s.libcats.NewQuery().
		Filter("Category = ", name).
		Run(&libcats)

	if err != nil {
		log.Printf("Error getting lib_categories: %s", err)
		c.AbortWithError(500, err)
		return
	}

	libs := []string{}

	for _, lc := range libcats {
		libs = append(libs, lc.Lib)
	}

	out := struct {
		models.Category
		Libs []string `json:"libs"`
	}{
		cat,
		libs,
	}

	c.JSON(200, out)
}
