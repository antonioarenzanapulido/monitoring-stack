package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests to the API",
		},
		[]string{"path", "method", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of API requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func main() {
	http.HandleFunc("/endpoint", withMetrics(errorRateHandler))
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func withMetrics(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture the response status code
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		handler.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		requestsTotal.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(rec.statusCode)).Inc()
		requestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
	}
}

func errorRateHandler(w http.ResponseWriter, r *http.Request) {
	// Get the error rate from the HTTP header
	errorRateHeader := r.Header.Get("X-Error-Rate")

	// Default to 0 if not set or invalid
	errorRate, err := strconv.Atoi(errorRateHeader)
	if err != nil || errorRate < 0 || errorRate > 100 {
		errorRate = 0
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 99
	randomNumber := rand.Intn(100)

	if randomNumber < errorRate {
		// Return HTTP 500 Internal Server Error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return HTTP 200 OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success!"))
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
