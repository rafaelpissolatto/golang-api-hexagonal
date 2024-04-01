package ports

import (
	"context"
	"golang-api-hexagonal/core/domain"
)

// IProductService product service interface
type IProductService interface {
	CreateProduct(ctx context.Context, request *domain.Product, username, traceID string) (*domain.ProductResponse, error)
	GetProduct(ctx context.Context, productID, traceID string) (*domain.ProductResponse, error)
}
