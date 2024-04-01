package domain

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// ProductEventName product kafka event name
const ProductEventName = "create.product.event"

// Product request product
type Product struct {
	Name        string `json:"name" validate:"required,not_blank,min=2,max=256"`
	Description string `json:"description" validate:"omitempty,min=2,max=256"`
	UnitType    string `json:"unitType" validate:"required,oneof=unit kilos grams liters box size"`
	Unit        string `json:"unit" validate:"required,not_blank,min=1,max=50"`
	Brand       string `json:"brand" validate:"required,not_blank,min=1,max=50"`
	Color       string `json:"color" validate:"required,not_blank,min=1,max=50"`
	Style       string `json:"style" validate:"required,not_blank,min=1,max=50"`
	Status      string `json:"status" validate:"required,oneof=available pending inactive"`
}

// ProductModel product database model
type ProductModel struct {
	bun.BaseModel `bun:"table:products" json:"-"`
	ID            string    `bun:"id,pk" json:"id"`
	Name          string    `bun:"name" json:"name"`
	Description   string    `bun:"description" json:"description"`
	UnitType      string    `bun:"unit_type" json:"unitType"`
	Unit          string    `bun:"unit" json:"unit"`
	Brand         string    `bun:"brand" json:"brand"`
	Color         string    `bun:"color" json:"color"`
	Style         string    `bun:"style" json:"style"`
	Status        string    `bun:"status" json:"status"`
	AuditUser     string    `bun:"audit_user" json:"auditUser"`
	CreationDate  time.Time `bun:"creation_date" json:"creationDate"`
	UpdateDate    time.Time `bun:"update_date" json:"updateDate"`
}

// ProductResponse product response
type ProductResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	UnitType     string    `json:"unitType"`
	Unit         string    `json:"unit"`
	Brand        string    `json:"brand"`
	Color        string    `json:"color"`
	Style        string    `json:"style"`
	Status       string    `json:"status"`
	AuditUser    string    `json:"auditUser"`
	CreationDate time.Time `json:"creationDate"`
	UpdateDate   time.Time `json:"updateDate"`
}

// FromProductToProductModel convert from Product Request to Product Database Model
func FromProductToProductModel(request *Product, auditUser string) *ProductModel {
	currentTime := time.Now()
	return &ProductModel{
		ID:           uuid.NewString(),
		Name:         request.Name,
		Description:  request.Description,
		UnitType:     request.UnitType,
		Unit:         request.Unit,
		Brand:        request.Brand,
		Color:        request.Color,
		Style:        request.Style,
		Status:       request.Status,
		AuditUser:    auditUser,
		CreationDate: currentTime,
		UpdateDate:   currentTime,
	}
}

// FromProductModelToProductResponse convert from Product database model to Product response
func FromProductModelToProductResponse(productModel *ProductModel) *ProductResponse {
	return &ProductResponse{
		ID:           productModel.ID,
		Name:         productModel.Name,
		Description:  productModel.Description,
		UnitType:     productModel.UnitType,
		Unit:         productModel.Unit,
		Brand:        productModel.Brand,
		Color:        productModel.Color,
		Style:        productModel.Style,
		Status:       productModel.Status,
		AuditUser:    productModel.AuditUser,
		CreationDate: productModel.CreationDate,
		UpdateDate:   productModel.UpdateDate,
	}
}
