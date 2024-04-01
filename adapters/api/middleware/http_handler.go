package middleware

import (
	"github.com/go-chi/chi/v5"
	md "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

// HttpMiddleware middleware for http requests
type HttpMiddleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewHttpHandlerMiddleware create a new http handler middleware
func NewHttpHandlerMiddleware(customMetricRegistry *CustomMetricRegistry) func(next http.Handler) http.Handler {
	// Create the http metric labels
	hm := &HttpMiddleware{
		reqs: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: httpServerSubsystem,
				Name:      metricReqsTotal,
				Help:      "How many http requests was processed by status code, method and path with patterns.",
			},
			[]string{code, method, path},
		),
		latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: httpServerSubsystem,
				Name:      metricLatency,
				Help:      "How long took to process the http requests by status code, method and path with patterns.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{code, method, path},
		),
	}
	// Register http request and latency metrics on prometheus
	customMetricRegistry.MustRegister(hm.reqs)
	customMetricRegistry.MustRegister(hm.latency)

	return hm.httpHandler
}

// httpHandler http handler route path
func (h HttpMiddleware) httpHandler(next http.Handler) http.Handler {
	// Fill the values in http metrics
	fc := func(writer http.ResponseWriter, reader *http.Request) {
		startTime := time.Now()
		wrapWriter := md.NewWrapResponseWriter(writer, reader.ProtoMajor)
		next.ServeHTTP(wrapWriter, reader)

		routectx := chi.RouteContext(reader.Context())
		if routectx != nil {
			h.reqs.WithLabelValues(strconv.Itoa(wrapWriter.Status()), reader.Method, routectx.RoutePattern()).Inc()
			h.latency.WithLabelValues(strconv.Itoa(wrapWriter.Status()), reader.Method, routectx.RoutePattern()).Observe(time.Since(startTime).Seconds())
		} else {
			h.reqs.WithLabelValues(strconv.Itoa(wrapWriter.Status()), reader.Method, reader.URL.Path).Inc()
			h.latency.WithLabelValues(strconv.Itoa(wrapWriter.Status()), reader.Method, reader.URL.Path).Observe(time.Since(startTime).Seconds())
		}
	}
	return http.HandlerFunc(fc)
}
