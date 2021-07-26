package middlewares

import (
	"context"
	"net/rpc"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/namo-io/kit/pkg/keys"
	"github.com/namo-io/kit/pkg/log/logger"
	"google.golang.org/grpc"
)

// Default middleware default set
func Default() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		Gzip(),
		RequestID(),
		CORS(),
	}
}

// RequestID setter RequestID(UUID) into route context
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New()

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, keys.RequestID, requestID.String())
		c.Request = c.Request.WithContext(ctx)
	}
}

// Logging log output for incoming connections
func Logging(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.WithContext(c.Request.Context())
		logger.Debugf("[%v] %s %s", c.Request.Method, c.Request.URL.RequestURI(), c.Request.URL.Query().Encode())
	}
}

// Gzip payload data compression
// if you want more detail, see https://ko.wikipedia.org/wiki/Gzip
func Gzip() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}

// CORS all allow origins (default)
func CORS() gin.HandlerFunc {
	return cors.Default()
}

func RPC(s *rpc.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.ServeHTTP(c.Writer, c.Request)
	}
}

func GRPC(s *grpc.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ProtoMajor == 2 && strings.Contains(c.ContentType(), "application/grpc") {
			s.ServeHTTP(c.Writer, c.Request)
		}
	}
}
