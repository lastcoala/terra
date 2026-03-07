package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/pkg/filter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormRepo(path string, nConn int) (*gorm.DB, error) {
	var err error
	client, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: client}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
	if err != nil {
		return nil, err
	}

	dbX, err := db.DB()
	if err != nil {
		return nil, err
	}
	dbX.SetMaxIdleConns(nConn)
	dbX.SetMaxOpenConns(nConn)
	dbX.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}

type ProductGormRepo struct {
	db *gorm.DB
}

func NewProductGormRepo(db *gorm.DB) IProductRepo {
	return &ProductGormRepo{db: db}
}

func (r ProductGormRepo) error(ctx context.Context, err error, method string, params ...any) error {
	errF := fmt.Errorf("ProductGormRepo.(%v)(%v) %w", method, params, err)
	slog.ErrorContext(ctx, errF.Error())
	return errF
}

func (r *ProductGormRepo) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	var product ProductGormModel
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return domain.Product{}, r.error(ctx, err, "GetProduct", id)
	}
	return product.ToDomain(), nil
}

func (r *ProductGormRepo) GetProducts(ctx context.Context, offset, limit int, filters ...filter.Filter) ([]domain.Product, error) {
	var products ProductGormModels
	query := r.db.WithContext(ctx)
	for _, f := range filters {
		query = query.Where(f.Attribute+" "+f.Operator+" ?", f.Value)
	}
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, r.error(ctx, err, "GetProducts")
	}
	return products.ToDomains(), nil
}

func (r *ProductGormRepo) InsertProduct(ctx context.Context, product domain.Product, createdBy string) (domain.Product, error) {
	model := NewProductGormModel(product)
	model.CreatedBy = createdBy
	model.UpdatedBy = createdBy

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Product{}, r.error(ctx, err, "InsertProduct", product.Name)
	}
	return model.ToDomain(), nil
}

func (r *ProductGormRepo) UpdateProduct(ctx context.Context, product domain.Product, updatedBy string) (domain.Product, error) {
	model := NewProductGormModel(product)
	model.UpdatedBy = updatedBy

	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.Product{}, r.error(ctx, err, "UpdateProduct", product.Name)
	}
	return model.ToDomain(), nil
}

func (r *ProductGormRepo) DeleteProduct(ctx context.Context, id int, deletedBy string) error {
	if err := r.db.WithContext(ctx).Model(&ProductGormModel{}).
		Where("id = ?", id).
		Updates(map[string]any{"is_deleted": true, "updated_by": deletedBy}).
		Error; err != nil {
		return r.error(ctx, err, "DeleteProduct", id)
	}
	return nil
}
