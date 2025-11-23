package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tenz-io/gokit/logger"
	"github.com/tenz-io/gokit/tracer"
)

// 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		ctx := c.Request.Context()
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = tracer.RequestIdFromCtx(ctx)
		}

		le := logger.FromContext(ctx).
			WithTracing(requestID).
			WithFields(logger.Fields{
				"path": path,
			})
		ctx = logger.WithLogger(ctx, le)

		te := logger.TrafficEntryFromContext(ctx).
			WithTracing(requestID).
			WithFields(logger.Fields{
				"path": path,
			})
		ctx = logger.WithTrafficEntry(ctx, te)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		logger.TrafficEntryFromContext(c.Request.Context()).
			WithFields(logger.Fields{
				"method":    c.Request.Method,
				"client_ip": c.ClientIP(),
			}).Data(&logger.Traffic{
			Typ:  logger.TrafficTypRecv,
			Cmd:  path,
			Cost: time.Now().Sub(start),
			Code: strconv.Itoa(c.Writer.Status()),
			Msg:  "request",
			Req:  raw,
		})
	}
}

// CORS 中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
