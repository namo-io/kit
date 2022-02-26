package middlewares

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/namo-io/kit/pkg/ctxkey"
	"github.com/namo-io/kit/pkg/log"
	"github.com/namo-io/kit/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

// Default middleware default set
func Default() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		CORS(),
		Gzip(),
		InjectRequestID(),
		InjectAuthorization(),
		Metrics(prometheus.DefaultRegisterer),
	}
}

// InjectRequestID inject request-id(UUID) into route context
// pass the requets id header, create if none exists
func InjectRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// read := x-request-id > request-id => new uuid()

		requestId := c.GetHeader("x-request-id")
		if requestId == "" {
			requestId = c.GetHeader(ctxkey.RequestId)
		}

		if requestId == "" {
			requestId = uuid.New().String()
		}

		// context
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ctxkey.RequestId, requestId)
		c.Request = c.Request.WithContext(ctx)
	}
}

// InjectAuthorization inject authorization into route context
func InjectAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader(ctxkey.Authorization)
		if authorization == "" {
			return
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ctxkey.Authorization, authorization)
		c.Request = c.Request.WithContext(ctx)
	}
}

// Logging log output for incoming connections
func Logging(log log.Log) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := log.WithContext(c.Request.Context())
		log.Debugf("[%v] %s %s", c.Request.Method, c.Request.URL.RequestURI(), c.Request.URL.Query().Encode())
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

func GRPC(s *grpc.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ProtoMajor == 2 && strings.HasPrefix("application/grpc", c.ContentType()) {
			s.ServeHTTP(c.Writer, c.Request)
		}
	}
}

// Metrics default metrics(http request total/latency ...)
func Metrics(register prometheus.Registerer) gin.HandlerFunc {
	hostname := util.GetHostname()

	ServerRequestStartAt := promauto.With(register).NewHistogramVec(prometheus.HistogramOpts{
		Name: "server_request_start_at",
	}, []string{"hostname", "method", "path", "status"})
	ServerRequestEndAt := promauto.With(register).NewHistogramVec(prometheus.HistogramOpts{
		Name: "server_request_end_at",
	}, []string{"hostname", "method", "path", "status"})

	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.RequestURI()
		status := fmt.Sprint(c.Writer.Status())

		t := float64(time.Now().UnixNano() / int64(time.Millisecond))
		ServerRequestStartAt.WithLabelValues(hostname, method, path, status).Observe(t)
		c.Next()

		t = float64(time.Now().UnixNano() / int64(time.Millisecond))
		ServerRequestEndAt.WithLabelValues(hostname, method, path, status).Observe(t)
	}
}
