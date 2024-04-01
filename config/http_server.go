package config

import (
	"go.uber.org/zap"
	"golang-api-hexagonal/adapters/api/router"
	"net/http"
	"time"
)

// StartHttpServer Config and start the http server
func StartHttpServer(log *zap.SugaredLogger, config ServerConfigurations, router *router.HTTPRouter) {
	log.Infof("Http server listening on port: %s", config.Port)

	// Config server
	server := &http.Server{
		Addr:         ":" + config.Port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      router.Router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start the http server: %v", err)
	}
}
