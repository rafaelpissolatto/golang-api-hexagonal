package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang-api-hexagonal/adapters/cache"
	"os"
	"path/filepath"
	"time"
)

// NewRedisCache creates a new cache redis connection
func NewRedisCache(log *zap.SugaredLogger, config RedisConfiguration) *cache.RedisCache {
	connection := connectToRedis(log, config)

	return &cache.RedisCache{
		Client: connection,
	}
}

func connectToRedis(log *zap.SugaredLogger, config RedisConfiguration) *redis.Client {
	var connection *redis.Client
	timeout := time.Duration(config.TimeOutInSeconds)

	if config.Localhost {
		connection = redis.NewClient(&redis.Options{
			Addr:         config.URL,
			DB:           config.DB,
			DialTimeout:  timeout * time.Second,
			ReadTimeout:  timeout * time.Second,
			WriteTimeout: timeout * time.Second,
			MaxRetries:   -1,
		})
	} else {
		cert, err1 := tls.LoadX509KeyPair(config.PublicKeyFile, config.PrivateKeyFile)
		if err1 != nil {
			log.Fatalf("Failed to load Redis Key Pairs: %v", err1)
		}

		caCert, err2 := os.ReadFile(filepath.Clean(config.CaCertFile))
		if err2 != nil {
			log.Fatalf("Failed to load Redis Cert: %v", err2)
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
			MinVersion:   tls.VersionTLS12,
		}

		connection = redis.NewClient(&redis.Options{
			Addr:         config.URL,
			Username:     config.User,
			Password:     config.Pass,
			DB:           config.DB,
			TLSConfig:    tlsConfig,
			DialTimeout:  timeout * time.Second,
			ReadTimeout:  timeout * time.Second,
			WriteTimeout: timeout * time.Second,
			MaxRetries:   -1,
		})
	}

	err3 := connection.Ping(context.Background()).Err()
	if err3 != nil {
		log.Fatalf("Failed to Ping Redis: %v", err3)
	}

	log.Infof("Redis connected!")
	return connection
}
