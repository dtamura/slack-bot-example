package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
)

func healthz(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "healthz")
	defer span.Finish()

	hostname, _ := os.Hostname()

	span.SetTag("hostname", hostname)

	// 成功時
	c.JSON(200, gin.H{
		"timestamp": time.Now(),
		"status":    "OK",
		"message":   "i'm healthy",
		"hostname":  hostname,
	})
}
