package products

import (
	"context"
	"golang-api-hexagonal/core/domain"
)

// ProductRepositoryMock product repository mock
type ProductRepositoryMock struct{}

var (
	CreateFunc              func(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error)
	ProductAlreadyExistFunc func(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error)
	GetProductByIdFunc      func(ctx context.Context, productID string) (*domain.ProductModel, error)
)

// Create is the repository mock for Create func
func (pr *ProductRepositoryMock) Create(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error) {
	return CreateFunc(ctx, model)
}

// ProductAlreadyExist is the repository mock for ProductAlreadyExist func
func (pr *ProductRepositoryMock) ProductAlreadyExist(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error) {
	return ProductAlreadyExistFunc(ctx, name, unitType, unit, brand, color, style)
}

// GetProductById is the repository mock for GetProductById func
func (pr *ProductRepositoryMock) GetProductById(ctx context.Context, productID string) (*domain.ProductModel, error) {
	return GetProductByIdFunc(ctx, productID)
}
