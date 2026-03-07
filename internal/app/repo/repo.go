package repo

import (
	"context"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/pkg/filter"
)

type IUserRepo interface {
	GetUser(ctx context.Context, id int) (domain.User, error)
	GetUsers(ctx context.Context, offset, limit int, filters ...filter.Filter) ([]domain.User, error)
	InsertUser(ctx context.Context, user domain.User, createdBy string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User, updatedBy string) (domain.User, error)
	DeleteUser(ctx context.Context, id int, deletedBy string) error
}
