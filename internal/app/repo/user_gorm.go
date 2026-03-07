package repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/pkg/filter"
	"gorm.io/gorm"
)

type UserGormRepo struct {
	db *gorm.DB
}

func NewUserGormRepo(db *gorm.DB) IUserRepo {
	return &UserGormRepo{db: db}
}

func (r UserGormRepo) error(ctx context.Context, err error, method string, params ...any) error {
	errF := fmt.Errorf("UserGormRepo.(%v)(%v) %w", method, params, err)
	slog.ErrorContext(ctx, errF.Error())
	return errF
}

func (r *UserGormRepo) GetUser(ctx context.Context, id int) (domain.User, error) {
	var user UserGormModel
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return domain.User{}, r.error(ctx, err, "GetUser", id)
	}
	return user.ToDomain(), nil
}

func (r *UserGormRepo) GetUsers(ctx context.Context, offset, limit int, filters ...filter.Filter) ([]domain.User, error) {
	var users UserGormModels
	query := r.db.WithContext(ctx).Order("id ASC")
	for _, filter := range filters {
		query = query.Where(filter.Attribute+" "+filter.Operator+" ?", filter.Value)
	}
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, r.error(ctx, err, "GetUsers")
	}
	return users.ToDomains(), nil
}

func (r *UserGormRepo) InsertUser(ctx context.Context, user domain.User, createdBy string) (domain.User, error) {
	model := NewUserGormModel(user)
	model.CreatedBy = createdBy
	model.UpdatedBy = createdBy

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.User{}, r.error(ctx, err, "InsertUser", user.Name)
	}
	return model.ToDomain(), nil
}

func (r *UserGormRepo) UpdateUser(ctx context.Context, user domain.User, updatedBy string) (domain.User, error) {
	model := NewUserGormModel(user)
	model.UpdatedBy = updatedBy

	query := r.db.WithContext(ctx)
	if user.Password == "" {
		query = query.Omit("password")
	}

	if err := query.Save(&model).Error; err != nil {
		return domain.User{}, r.error(ctx, err, "UpdateUser", user.Name)
	}
	return model.ToDomain(), nil
}

func (r *UserGormRepo) DeleteUser(ctx context.Context, id int, deletedBy string) error {
	if err := r.db.WithContext(ctx).Model(&UserGormModel{}).
		Where("id = ?", id).
		Updates(map[string]any{"is_deleted": true, "updated_by": deletedBy}).
		Error; err != nil {
		return r.error(ctx, err, "DeleteUser", id)
	}
	return nil
}
