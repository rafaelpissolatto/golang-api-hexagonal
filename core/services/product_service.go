package services

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"golang-api-hexagonal/adapters/cache"
	"golang-api-hexagonal/adapters/custom_error"
	"golang-api-hexagonal/config"
	"golang-api-hexagonal/core/domain"
	"golang-api-hexagonal/core/ports"
	"net/http"
)

// ProductService product service
type ProductService struct {
	log               *zap.SugaredLogger
	productRepository ports.IRepository
	redis             ports.IRedis
	message           ports.IMessage
	messageConfig     config.KafkaConfiguration
}

// NewProductService create new product service
func NewProductService(log *zap.SugaredLogger, productRepository ports.IRepository, redis ports.IRedis, message ports.IMessage,
	messageConfig config.KafkaConfiguration) *ProductService {
	return &ProductService{
		log:               log,
		productRepository: productRepository,
		redis:             redis,
		message:           message,
		messageConfig:     messageConfig,
	}
}

// CreateProduct service to create the product
func (ps *ProductService) CreateProduct(ctx context.Context, request *domain.Product, username, traceID string) (*domain.ProductResponse, error) {
	exist, err := ps.productRepository.ProductAlreadyExist(ctx,
		request.Name, request.UnitType, request.Unit, request.Brand, request.Color, request.Style)
	if err != nil {
		ps.log.With("traceId", traceID).Errorf("Internal server error: %v", err)
		return nil, custom_error.New(http.StatusInternalServerError, "internal server error")
	} else if exist {
		ps.log.With("traceId", traceID).Errorf("Product already exist")
		return nil, custom_error.New(http.StatusConflict, "already exist")
	}

	productModel := domain.FromProductToProductModel(request, username)

	_, err = ps.productRepository.Create(ctx, productModel)
	if err != nil {
		ps.log.With("traceId", traceID).Errorf("Internal server error: %v", err)
		return nil, custom_error.New(http.StatusInternalServerError, "internal server error")
	}

	data, errMarshall := json.Marshal(productModel)
	if errMarshall != nil {
		ps.log.With("traceId", traceID).Errorf("Internal error to marshal the payload: %v", errMarshall)
	} else {
		errCache := ps.redis.Set(ctx, productModel.ID, data, cache.KeyCacheDuration).Err()
		if errCache != nil {
			ps.log.With("traceId", traceID).Errorf("Internal error to save in cache: %v", errCache)
		}
	}

	ps.message.ProduceMessage(ps.messageConfig.Producer.ProductTopic, string(data), domain.ProductEventName, traceID)

	ps.log.With("traceId", traceID).Infof("The productID %s was created with success", productModel.ID)
	return &domain.ProductResponse{ID: productModel.ID}, nil
}

// GetProduct get the product by id
func (ps *ProductService) GetProduct(ctx context.Context, productID, traceID string) (*domain.ProductResponse, error) {
	var product *domain.ProductModel

	payloadBytes, errCache := ps.redis.Get(ctx, productID).Bytes()
	if errCache != nil {
		ps.log.With("traceId", traceID).Infof("productID not found in cache: %v", errCache)
	} else {
		errMar := json.Unmarshal(payloadBytes, &product)
		if errMar != nil {
			ps.log.With("traceId", traceID).Errorf("Internal error unmarshal the payload: %v", errMar)
		} else {
			ps.log.With("traceId", traceID).Infof("The productID %s was found with success in cache", product.ID)
			return domain.FromProductModelToProductResponse(product), nil
		}
	}

	productModel, err := ps.productRepository.GetProductById(ctx, productID)
	if err != nil {
		ps.log.With("traceId", traceID).Errorf("Internal server error to get the product: %v", err)
		return nil, custom_error.New(http.StatusInternalServerError, "internal server error")
	}
	if productModel == nil {
		ps.log.With("traceId", traceID).Errorf("Product not found")
		return nil, custom_error.New(http.StatusNotFound, "not found")
	}

	ps.log.With("traceId", traceID).Infof("The productID %s was found with success", productModel.ID)
	return domain.FromProductModelToProductResponse(productModel), nil
}
