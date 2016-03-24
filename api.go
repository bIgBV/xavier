package main

import (
	"github.com/gin-gonic/gin"
)

func serveMetrics() {
	router := gin.Default()
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}
