package service

import (
	"context"

	"github.com/lastcoala/terra/internal/app/domain"
)

type IUserService interface {
	InsertUser(ctx context.Context, user domain.User, createdBy string) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	GetUsers(ctx context.Context, offset, limit int) ([]domain.User, error)
	UpdateUser(ctx context.Context, user domain.User, updatedBy string) (domain.User, error)
	ChangePassword(ctx context.Context, id int, newPassword string, updatedBy string) (domain.User, error)
	DeleteUser(ctx context.Context, id int, deletedBy string) error
}
