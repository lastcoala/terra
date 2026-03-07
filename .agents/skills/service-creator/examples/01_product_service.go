package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/internal/app/repo"
)

type ProductService struct {
	repo repo.IProductRepo
}

func NewProductService(repo repo.IProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s ProductService) error(ctx context.Context, err error, method string, params ...any) error {
	errF := fmt.Errorf("ProductService.(%v)(%v) %w", method, params, err)
	slog.ErrorContext(ctx, errF.Error())
	return errF
}

func (s ProductService) InsertProduct(ctx context.Context, product domain.Product, createdBy string) (domain.Product, error) {
	if err := product.Validate(); err != nil {
		return domain.Product{}, s.error(ctx, err, "InsertProduct", product.Name)
	}

	created, err := s.repo.InsertProduct(ctx, product, createdBy)
	if err != nil {
		return domain.Product{}, s.error(ctx, err, "InsertProduct", product.Name)
	}
	return created, nil
}

func (s ProductService) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return domain.Product{}, s.error(ctx, err, "GetProduct", id)
	}
	return product, nil
}

func (s ProductService) GetProducts(ctx context.Context, offset, limit int) ([]domain.Product, error) {
	products, err := s.repo.GetProducts(ctx, offset, limit)
	if err != nil {
		return nil, s.error(ctx, err, "GetProducts", offset, limit)
	}
	return products, nil
}

func (s ProductService) UpdateProduct(ctx context.Context, product domain.Product, updatedBy string) (domain.Product, error) {
	// Fetch the existing record to preserve any fields the caller should not overwrite.
	existing, err := s.repo.GetProduct(ctx, product.Id)
	if err != nil {
		return domain.Product{}, s.error(ctx, err, "UpdateProduct", product.Id)
	}

	// Example: preserve a sensitive field from the stored record.
	_ = existing // replace with real field preservation as needed

	if err := product.Validate(); err != nil {
		return domain.Product{}, s.error(ctx, err, "UpdateProduct", product.Id)
	}

	updated, err := s.repo.UpdateProduct(ctx, product, updatedBy)
	if err != nil {
		return domain.Product{}, s.error(ctx, err, "UpdateProduct", product.Id)
	}
	return updated, nil
}

func (s ProductService) DeleteProduct(ctx context.Context, id int, deletedBy string) error {
	if err := s.repo.DeleteProduct(ctx, id, deletedBy); err != nil {
		return s.error(ctx, err, "DeleteProduct", id)
	}
	return nil
}
