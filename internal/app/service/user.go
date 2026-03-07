package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/internal/app/repo"
)

type UserService struct {
	repo repo.IUserRepo
}

func NewUserService(repo repo.IUserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s UserService) error(ctx context.Context, err error, method string, params ...any) error {
	errF := fmt.Errorf("UserService.(%v)(%v) %w", method, params, err)
	slog.ErrorContext(ctx, errF.Error())
	return errF
}

func (s UserService) InsertUser(ctx context.Context, user domain.User, createdBy string) (domain.User, error) {
	err := user.Validate(true)
	if err != nil {
		return domain.User{}, s.error(ctx, err, "InsertUser", user.Name)
	}

	err = user.EncryptPassword()
	if err != nil {
		return domain.User{}, s.error(ctx, err, "InsertUser", user.Name)
	}

	created, err := s.repo.InsertUser(ctx, user, createdBy)
	if err != nil {
		return domain.User{}, s.error(ctx, err, "InsertUser", user.Name)
	}

	return created, nil
}

func (s UserService) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, s.error(ctx, err, "GetUser", id)
	}
	return user, nil
}

func (s UserService) GetUsers(ctx context.Context, offset, limit int) ([]domain.User, error) {
	users, err := s.repo.GetUsers(ctx, offset, limit)
	if err != nil {
		return nil, s.error(ctx, err, "GetUsers", offset, limit)
	}
	return users, nil
}

func (s UserService) UpdateUser(ctx context.Context, user domain.User, updatedBy string) (domain.User, error) {
	existingUser, err := s.repo.GetUser(ctx, user.Id)
	if err != nil {
		return domain.User{}, s.error(ctx, err, "UpdateUser", user.Id)
	}

	user.Password = existingUser.Password

	err = user.Validate(false)
	if err != nil {
		return domain.User{}, s.error(ctx, err, "UpdateUser", user.Id)
	}

	updated, err := s.repo.UpdateUser(ctx, user, updatedBy)
	if err != nil {
		return domain.User{}, s.error(ctx, err, "UpdateUser", user.Id)
	}
	return updated, nil
}

func (s UserService) ChangePassword(ctx context.Context, id int, newPassword string, updatedBy string) (domain.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, s.error(ctx, err, "ChangePassword", id)
	}

	user.Password = newPassword
	err = user.ValidatePassword()
	if err != nil {
		return domain.User{}, s.error(ctx, err, "ChangePassword", id)
	}

	err = user.EncryptPassword()
	if err != nil {
		return domain.User{}, s.error(ctx, err, "ChangePassword", id)
	}

	updated, err := s.repo.UpdateUser(ctx, user, updatedBy)
	if err != nil {
		return domain.User{}, s.error(ctx, err, "ChangePassword", id)
	}
	return updated, nil
}

func (s UserService) DeleteUser(ctx context.Context, id int, deletedBy string) error {
	err := s.repo.DeleteUser(ctx, id, deletedBy)
	if err != nil {
		return s.error(ctx, err, "DeleteUser", id)
	}
	return nil
}
