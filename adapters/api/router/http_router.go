package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.elastic.co/apm/module/apmchiv5/v2"
	middleware3 "golang-api-hexagonal/adapters/api/middleware"
	"time"
)

// HTTPRouter http routers
type HTTPRouter struct {
	Router *chi.Mux
}

// NewHTTPRouter create new http router
func NewHTTPRouter(prometheusMetricRegistry *middleware3.CustomMetricRegistry) *HTTPRouter {
	router := chi.NewRouter()

	// Adding some middlewares ready
	router.Use(middleware.AllowContentType("application/json"))
	// Enable Elastic APM chiv5 Middleware
	router.Use(apmchiv5.Middleware())
	// Timeout is a middleware that cancels ctx after a given timeout and return http status error 504.
	router.Use(middleware.Timeout(60 * time.Second))
	// RequestID is a middleware that injects a request ID into the context of each request.
	router.Use(middleware.RequestID)
	// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
	router.Use(middleware.RealIP)
	// Recover is a middleware that recovers from panics, logs the panic, and returns an HTTP 500.
	router.Use(middleware.Recoverer)
	// Use a http middleware to pattern request for prometheus
	router.Use(middleware3.NewHttpHandlerMiddleware(prometheusMetricRegistry))

	return &HTTPRouter{
		Router: router,
	}
}
