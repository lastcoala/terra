package repo

import "github.com/lastcoala/terra/internal/app/domain"

type ProductGormModel struct {
	BaseModel
	Name        string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
}

func NewProductGormModel(p domain.Product) ProductGormModel {
	return ProductGormModel{
		BaseModel:   BaseModel{Id: p.Id},
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}
}

func (m ProductGormModel) ToDomain() domain.Product {
	return domain.Product{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
	}
}

func (m ProductGormModel) TableName() string {
	return "products"
}

type ProductGormModels []ProductGormModel

func (m ProductGormModels) ToDomains() []domain.Product {
	var products []domain.Product
	for _, p := range m {
		products = append(products, p.ToDomain())
	}
	return products
}

func (m ProductGormModels) TableName() string {
	return "products"
}
