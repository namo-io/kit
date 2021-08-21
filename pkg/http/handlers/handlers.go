package handlers

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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

func Graphql(es graphql.ExecutableSchema) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c.Request.Context())
		s := handler.NewDefaultServer(es)
		s.Use(handler.ResponseFunc(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			defer cancel()
			return next(ctx)
		}))
		s.ServeHTTP(c.Writer, c.Request)

		<-ctx.Done()
	}
}

func GraphqlPlayground() gin.HandlerFunc {
	return func(c *gin.Context) {
		playground.Handler("GraphqlQL PlayGround", c.Request.URL.Path).ServeHTTP(c.Writer, c.Request)
	}
}
