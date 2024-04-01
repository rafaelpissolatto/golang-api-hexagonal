package main

import (
	"context"
	"golang-api-hexagonal/adapters/api/controller"
	middleware2 "golang-api-hexagonal/adapters/api/middleware"
	"golang-api-hexagonal/adapters/api/router"
	"golang-api-hexagonal/adapters/kafka"
	"golang-api-hexagonal/adapters/opa"
	"golang-api-hexagonal/adapters/repository/products"
	"golang-api-hexagonal/config"
	"golang-api-hexagonal/core/services"

	"github.com/go-playground/validator/v10"
)

func main() {
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	configs := config.LoadConfigFile(logger)

	database := config.NewDatabaseConnection(logger, configs.DB)
	defer config.CloseDatabaseConnection(database)

	// Opa Policies
	policies := opa.NewPolicyService(configs.Policies.Path, logger)

	// Redis
	redisCache := config.NewRedisCache(logger, configs.Redis)

	// Repositories
	productsRepository := products.NewProductRepository(database)

	// Start Kafka Producer and Consumer with a new context
	ctx := context.Background()
	kafka.CreateKafkaTopics(logger, configs.Kafka, ctx, config.NewKafkaConfigMap(logger, configs.Kafka, config.Topic))
	producer := kafka.NewKafkaProducer(logger, config.NewKafkaConfigMap(logger, configs.Kafka, config.Producer))
	defer producer.Close()
	consumer := kafka.NewKafkaConsumer(logger, config.NewKafkaConfigMap(logger, configs.Kafka, config.Consumer))
	defer kafka.CloseConsumer(consumer)
	go kafka.ConsumeMessages(logger, configs.Kafka, consumer)

	// Config Domain Services
	productService := services.NewProductService(logger, productsRepository, redisCache, producer, configs.Kafka)
	authService := services.NewAuthService(logger, configs.Oauth)

	// Metrics
	prometheusMetrics := middleware2.NewPrometheusMiddleware(configs.Service.Name)
	jwtHandler := middleware2.NewJWTHandler(logger, authService)

	// Config Http Routers and Controllers
	route := router.NewHTTPRouter(prometheusMetrics)
	valid := validator.New()
	controller.NewHealthCheckController(route, prometheusMetrics)
	controller.NewAuthController(route, logger, valid, authService)
	controller.NewProductController(route, logger, valid, prometheusMetrics, productService, jwtHandler, policies)

	config.StartHttpServer(logger, configs.Server, route)
}
