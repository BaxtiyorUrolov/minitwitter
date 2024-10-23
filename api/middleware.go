package api

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

func requestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)
		log.Println(string(body))
		log.Println(c.Request.Header)
		c.Next()
	}
}
