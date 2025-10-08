// main.go
package main

import (
	"math/rand"
	"net/http"
	"os"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fintech_requests_total",
			Help: "Total number of requests",
		},
		[]string{"endpoint", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "fintech_request_duration_seconds",
			Help:    "Request duration distribution",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestDuration)
}

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		start := time.Now() //new changes
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
		duration := time.Since(start).Seconds()                      //new changes
		requestDuration.WithLabelValues("/health").Observe(duration) //new changes
		totalRequests.WithLabelValues("/health", "200").Inc()        //new changes
	})

	// Simulate transaction
	r.POST("/transaction", func(c *gin.Context) {
		// Example fintech processing
		start := time.Now()
		amount := c.Query("amount")
		if amount == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "amount required"})
			duration := time.Since(start).Seconds()
			requestDuration.WithLabelValues("/transaction").Observe(duration)
			totalRequests.WithLabelValues("/transaction", "400").Inc()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaction processed", "amount": amount})
		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues("/transaction").Observe(duration)
		totalRequests.WithLabelValues("/transaction", "200").Inc()
	})

	//FailureRoute(r) // <-- add simulated failure

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

	r.GET("/metrics", gin.WrapH(promhttp.Handler())) //new changes

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	r.Run(":" + port)
}
