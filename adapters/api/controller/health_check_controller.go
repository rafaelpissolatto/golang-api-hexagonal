package controller

import (
	"golang-api-hexagonal/adapters/api/dto"
	"golang-api-hexagonal/adapters/api/middleware"
	"golang-api-hexagonal/adapters/api/router"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewHealthCheckController create a new health check http controller
func NewHealthCheckController(httpRouter *router.HTTPRouter, prometheusRegistry *middleware.CustomMetricRegistry) {
	// health check endpoints for kubernetes
	httpRouter.Router.Get("/health/live", handleLivelinessCheck)
	httpRouter.Router.Get("/health/ready", handleReadinessCheck)
	// prometheus metrics endpoint
	httpRouter.Router.Get("/metrics", promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{}).ServeHTTP)
}

func handleLivelinessCheck(writer http.ResponseWriter, reader *http.Request) {
	dto.RenderResponse(reader.Context(), writer, http.StatusOK, dto.DefaultResponse(http.StatusText(http.StatusOK), ""))
}

func handleReadinessCheck(writer http.ResponseWriter, reader *http.Request) {
	dto.RenderResponse(reader.Context(), writer, http.StatusOK, dto.DefaultResponse(http.StatusText(http.StatusOK), ""))
}
