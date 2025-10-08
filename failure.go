package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func FailureRoute(r *gin.Engine) {
	r.GET("/simulate", func(c *gin.Context) {
		start := time.Now()
		if rand.Intn(100) < 80 { // 80% failure
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
			duration := time.Since(start).Seconds()
			requestDuration.WithLabelValues("/simulate").Observe(duration)
			totalRequests.WithLabelValues("/simulate", "500").Inc()
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Transaction success"})
		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues("/simulate").Observe(duration)
		totalRequests.WithLabelValues("/simulate", "200").Inc()
	})
}
