package products

import (
	"context"
	"errors"
	"github.com/uptrace/bun"
	"golang-api-hexagonal/core/domain"
	"strings"
	"sync"
)

// ProductRepository repository implementation for products
type ProductRepository struct {
	db         bun.IDB
	lockSelect sync.RWMutex
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(db bun.IDB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Create a new product
func (repo *ProductRepository) Create(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error) {
	resp, err := repo.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}
	affectedRows, err := resp.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affectedRows == 0 {
		return nil, errors.New("no rows inserted")
	}
	return model, nil
}

// ProductAlreadyExist product already exist?
func (repo *ProductRepository) ProductAlreadyExist(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error) {
	var product domain.ProductModel
	repo.lockSelect.RLock()

	err := repo.db.NewSelect().
		Model((*domain.ProductModel)(nil)).
		Where("name = ?", name).
		Where("unit_type = ?", unitType).
		Where("unit = ?", unit).
		Where("brand = ?", brand).
		Where("color = ?", color).
		Where("style = ?", style).
		Scan(ctx, &product)

	repo.lockSelect.RUnlock()
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetProductById get the product by id
func (repo *ProductRepository) GetProductById(ctx context.Context, productID string) (*domain.ProductModel, error) {
	var product domain.ProductModel
	repo.lockSelect.RLock()

	err := repo.db.NewSelect().
		Model((*domain.ProductModel)(nil)).
		Where("id = ?", productID).
		Scan(ctx, &product)

	repo.lockSelect.RUnlock()
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}
