package services

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"golang-api-hexagonal/adapters/cache"
	"golang-api-hexagonal/adapters/kafka"
	"golang-api-hexagonal/adapters/repository/products"
	"golang-api-hexagonal/config"
	"golang-api-hexagonal/core/domain"
	"testing"
	"time"
)

var log = config.NewLogger()
var defaultContext = context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
var username = "user_test"
var traceID = "64675856476878"
var product = &domain.Product{Name: "product_test"}

// TestCreateProductThatAlreadyExistError for test CreateProduct
func TestCreateProductThatAlreadyExistError(t *testing.T) {
	service := NewProductService(log, &products.ProductRepositoryMock{}, &cache.RedisCacheMock{}, &kafka.MessageProducerMock{}, config.KafkaConfiguration{})

	products.ProductAlreadyExistFunc = func(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error) {
		return true, nil
	}

	_, err := service.CreateProduct(defaultContext, product, username, traceID)
	assert.Equal(t, err.Error(), "already exist")
}

// TestCreateProductWithInternalServerErrorToFound for test CreateProduct
func TestCreateProductWithInternalServerErrorToFound(t *testing.T) {
	service := NewProductService(log, &products.ProductRepositoryMock{}, &cache.RedisCacheMock{}, &kafka.MessageProducerMock{}, config.KafkaConfiguration{})

	products.ProductAlreadyExistFunc = func(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error) {
		return false, errors.New("internal query error")
	}

	_, err := service.CreateProduct(defaultContext, product, username, traceID)
	assert.Equal(t, err.Error(), "internal server error")
}

// TestCreateProductWithInternalServerErrorToSave for test CreateProduct
func TestCreateProductWithInternalServerErrorToSave(t *testing.T) {
	service := NewProductService(log, &products.ProductRepositoryMock{}, &cache.RedisCacheMock{}, &kafka.MessageProducerMock{}, config.KafkaConfiguration{})

	products.ProductAlreadyExistFunc = func(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error) {
		return false, nil
	}

	products.CreateFunc = func(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error) {
		return nil, errors.New("internal error to save")
	}

	_, err := service.CreateProduct(defaultContext, product, username, traceID)
	assert.Equal(t, err.Error(), "internal server error")
}

// TestCreateProductWithSuccessButFailToCache for test CreateProduct
func TestCreateProductWithSuccessButFailToCache(t *testing.T) {
	service := NewProductService(log, &products.ProductRepositoryMock{}, &cache.RedisCacheMock{}, &kafka.MessageProducerMock{}, config.KafkaConfiguration{})

	products.ProductAlreadyExistFunc = func(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error) {
		return false, nil
	}

	products.CreateFunc = func(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error) {
		return nil, nil
	}

	var redisResponse = &redis.StatusCmd{}
	redisResponse.SetErr(errors.New("internal error"))
	cache.SetFunc = func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
		return redisResponse
	}

	kafka.ProduceMessageFunc = func(topicName, value, eventName, traceID string) {}

	productResponse, _ := service.CreateProduct(defaultContext, product, username, traceID)
	assert.NotEmpty(t, productResponse.ID)
}

// TestCreateProductWithSuccess for test CreateProduct
func TestCreateProductWithSuccess(t *testing.T) {
	service := NewProductService(log, &products.ProductRepositoryMock{}, &cache.RedisCacheMock{}, &kafka.MessageProducerMock{}, config.KafkaConfiguration{})

	products.ProductAlreadyExistFunc = func(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error) {
		return false, nil
	}

	products.CreateFunc = func(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error) {
		return nil, nil
	}

	var redisResponse = &redis.StatusCmd{}
	cache.SetFunc = func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
		return redisResponse
	}

	kafka.ProduceMessageFunc = func(topicName, value, eventName, traceID string) {}

	productResponse, _ := service.CreateProduct(defaultContext, product, username, traceID)
	assert.NotEmpty(t, productResponse.ID)
}
