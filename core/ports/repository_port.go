package ports

import (
	"context"
	"golang-api-hexagonal/core/domain"
)

// IRepository repository interface
type IRepository interface {
	Create(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error)
	ProductAlreadyExist(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error)
	GetProductById(ctx context.Context, productID string) (*domain.ProductModel, error)
}
