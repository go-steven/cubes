package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"res": "Cubes web server works.",
	})
}
