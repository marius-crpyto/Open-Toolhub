package server

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/marius-crpyto/toolhub/logger"
)

type GinLogger struct {
	Logger *logger.Logger
}

func (g GinLogger) Write(p []byte) (n int, err error) {
	g.Logger.Info("HTTP request",
		zap.String("log", strings.TrimSpace(string(p))),
	)
	return len(p), nil
}

func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		ctx.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		ctx.Next()
	}
}

func CORS(origins, allowHeaders []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: origins, //[]string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: append([]string{"Origin", "Content-Length", "Content-Type", "Authorization"}, allowHeaders...),
		MaxAge:       12 * time.Hour,
	})
}

func CORSWithOptions(origins, allowHeaders, exposeHeaders []string, allowCredentials bool) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     append([]string{"Origin", "Content-Length", "Content-Type", "Authorization"}, allowHeaders...),
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	})
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Request.Header.Get("X-Request-ID")
		if rid == "" {
			b := make([]byte, 16)
			_, _ = rand.Read(b)
			rid = hex.EncodeToString(b)
		}
		c.Set("request_id", rid)
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Next()
	}
}

func AccessLog(l *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		ip := c.ClientIP()
		method := c.Request.Method
		ua := c.Request.UserAgent()
		rid := c.GetString("request_id")
		in := c.Request.ContentLength
		if in < 0 {
			in = 0
		}
		out := c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		l.Info("access",
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.Duration("latency", latency),
			zap.String("client_ip", ip),
			zap.String("user_agent", ua),
			zap.Int64("bytes_in", in),
			zap.Int("bytes_out", out),
			zap.String("request_id", rid),
		)
	}
}
