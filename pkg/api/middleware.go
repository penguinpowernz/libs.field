package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	c.Set("limit", uint(perp))
	c.Set("offset", uint((page-1)*perp))
}
