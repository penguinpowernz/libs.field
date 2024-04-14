package main

import (
	"log"

	"github.com/albrow/zoom"
	"github.com/gin-gonic/gin"
	API "github.com/penguinpowernz/libs.fieid/pkg/api"
	"github.com/penguinpowernz/libs.fieid/pkg/models"
	"github.com/penguinpowernz/libs.fieid/pkg/util"
)

func main() {
	pool := zoom.NewPool("redis:6379")
	defer pool.Close()

	opts := zoom.CollectionOptions{
		FallbackMarshalerUnmarshaler: util.FallbackMarshaler{},
		Index:                        true,
	}

	libs, err := pool.NewCollectionWithOptions(&models.Lib{}, opts)
	cats, err := pool.NewCollectionWithOptions(&models.Category{}, opts)
	libcats, err := pool.NewCollectionWithOptions(&models.LibCategory{}, opts)

	if err != nil {
		log.Fatal(err)
	}

	api := gin.Default()
	api.LoadHTMLGlob("/var/www/html/*")
	svr := API.NewServer(libs, cats, libcats)
	svr.SetupRoutes(api)

	if err := api.Run("0.0.0.0:80"); err != nil {
		log.Fatal(err)
	}

}
