package api

import (
	"github.com/albrow/zoom"
	"github.com/gin-gonic/gin"
)

type Server struct {
	libs    *zoom.Collection
	cats    *zoom.Collection
	libcats *zoom.Collection
}

func NewServer(libs, cats, libcats *zoom.Collection) *Server {
	return &Server{
		libs:    libs,
		cats:    cats,
		libcats: libcats,
	}
}

func (svr *Server) SetupRoutes(r gin.IRouter) {

	r.GET("/", paginate, svr.index)

	api := r.Group("/v1")
	api.GET("/topics")
	api.GET("/topic/:name")
	api.GET("/categories", svr.getCategories)
	api.GET("/category/:name", svr.getCategory)
	api.GET("/libs", paginate, svr.getLibs)
	api.GET("/lib/:owner/:name", svr.getLib)
	api.GET("/libs/popular", paginate, svr.getLibsPopular)
	api.GET("/libs/growing", paginate, svr.getLibsGrowing)
	api.GET("/libs/recent")
	api.GET("/libs/active")
}
