package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namo-io/kit/pkg/version"
)

func OK() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func Version() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.ContentType() {
		case gin.MIMEJSON:
			c.JSON(http.StatusOK, version.String())
		default:
			c.String(http.StatusOK, version.Info().String())
		}
	}
}
