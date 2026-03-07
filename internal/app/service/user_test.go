package service

import (
	"errors"
	"testing"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/internal/mocks"
	"github.com/lastcoala/terra/pkg/ctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var testCtx = ctx.NewCtx("1234567890")

// newService wires a fresh MockIUserRepo into a UserService.
func newService(t *testing.T) (*UserService, *mocks.MockIUserRepo) {
	t.Helper()
	mockRepo := mocks.NewMockIUserRepo(t)
	svc := NewUserService(mockRepo)
	return svc, mockRepo
}

// validUser returns a domain.User that satisfies all validation rules.
func validUser() domain.User {
	return domain.User{
		Id:       1,
		Name:     "Alice",
		Gender:   "female",
		Email:    "alice@example.com",
		Password: "Secret123",
	}
}

// ---------------------------------------------------------------------------
// InsertUser
// ---------------------------------------------------------------------------

func TestUserService_InsertUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newService(t)

		user := validUser()
		user.Id = 0

		expected := domain.User{Id: 1, Name: user.Name, Email: user.Email, Gender: user.Gender}

		// The service bcrypt-encrypts the password before calling the repo, so we
		// match on any user whose hash verifies against the original plain-text.
		mockRepo.EXPECT().
			InsertUser(testCtx, mock.MatchedBy(func(u domain.User) bool {
				return u.Email == user.Email &&
					bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("Secret123")) == nil
			}), "sysadmin").
			Return(expected, nil)

		created, err := svc.InsertUser(testCtx, user, "sysadmin")
		require.NoError(t, err)
		assert.Equal(t, 1, created.Id)
		assert.Equal(t, user.Email, created.Email)
	})

	t.Run("invalid_email_skips_repo", func(t *testing.T) {
		svc, _ := newService(t)

		user := validUser()
		user.Email = "not-an-email"

		_, err := svc.InsertUser(testCtx, user, "sysadmin")
		assert.Error(t, err)
		// No EXPECT() set → AssertExpectations catches any unexpected repo call.
	})

	t.Run("invalid_gender_skips_repo", func(t *testing.T) {
		svc, _ := newService(t)

		user := validUser()
		user.Gender = "other"

		_, err := svc.InsertUser(testCtx, user, "sysadmin")
		assert.Error(t, err)
	})

	t.Run("invalid_password_skips_repo", func(t *testing.T) {
		svc, _ := newService(t)

		user := validUser()
		user.Password = "weak"

		_, err := svc.InsertUser(testCtx, user, "sysadmin")
		assert.Error(t, err)
	})

	t.Run("repo_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newService(t)

		user := validUser()
		user.Id = 0

		mockRepo.EXPECT().
			InsertUser(testCtx, mock.Anything, "sysadmin").
			Return(domain.User{}, errors.New("db error"))

		_, err := svc.InsertUser(testCtx, user, "sysadmin")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// GetUser
// ---------------------------------------------------------------------------

func TestUserService_GetUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newService(t)

		expected := validUser()
		mockRepo.EXPECT().GetUser(testCtx, 1).Return(expected, nil)

		got, err := svc.GetUser(testCtx, 1)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not_found_propagates_error", func(t *testing.T) {
		svc, mockRepo := newService(t)

		mockRepo.EXPECT().GetUser(testCtx, 999).Return(domain.User{}, errors.New("record not found"))

		_, err := svc.GetUser(testCtx, 999)
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// GetUsers
// ---------------------------------------------------------------------------

func TestUserService_GetUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newService(t)

		expected := []domain.User{validUser()}
		mockRepo.EXPECT().GetUsers(testCtx, 0, 10).Return(expected, nil)

		got, err := svc.GetUsers(testCtx, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("empty_result", func(t *testing.T) {
		svc, mockRepo := newService(t)

		mockRepo.EXPECT().GetUsers(testCtx, 0, 10).Return([]domain.User{}, nil)

		got, err := svc.GetUsers(testCtx, 0, 10)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("repo_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newService(t)

		mockRepo.EXPECT().GetUsers(testCtx, 0, 10).Return(nil, errors.New("db error"))

		_, err := svc.GetUsers(testCtx, 0, 10)
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// UpdateUser
// ---------------------------------------------------------------------------

func TestUserService_UpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newService(t)

		existing := validUser()
		existing.Password = "StoredHash"

		input := existing
		input.Name = "AliceUpdated"

		expected := input

		// Service fills in the password from the stored record before calling repo.
		mockRepo.EXPECT().GetUser(testCtx, existing.Id).Return(existing, nil)
		mockRepo.EXPECT().
			UpdateUser(testCtx, mock.MatchedBy(func(u domain.User) bool {
				return u.Id == existing.Id &&
					u.Password == existing.Password &&
					u.Name == "AliceUpdated"
			}), "admin").
			Return(expected, nil)

		got, err := svc.UpdateUser(testCtx, input, "admin")
		require.NoError(t, err)
		assert.Equal(t, "AliceUpdated", got.Name)
	})

	t.Run("user_not_found", func(t *testing.T) {
		svc, mockRepo := newService(t)

		mockRepo.EXPECT().GetUser(testCtx, 99).Return(domain.User{}, errors.New("not found"))

		user := validUser()
		user.Id = 99
		_, err := svc.UpdateUser(testCtx, user, "admin")
		assert.Error(t, err)
	})

	t.Run("invalid_email_skips_repo_update", func(t *testing.T) {
		svc, mockRepo := newService(t)

		existing := validUser()
		mockRepo.EXPECT().GetUser(testCtx, existing.Id).Return(existing, nil)

		bad := existing
		bad.Email = "bad-email"
		_, err := svc.UpdateUser(testCtx, bad, "admin")
		assert.Error(t, err)
	})

	t.Run("repo_update_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newService(t)

		existing := validUser()
		mockRepo.EXPECT().GetUser(testCtx, existing.Id).Return(existing, nil)
		mockRepo.EXPECT().UpdateUser(testCtx, mock.Anything, "admin").Return(domain.User{}, errors.New("db error"))

		_, err := svc.UpdateUser(testCtx, existing, "admin")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// ChangePassword
// ---------------------------------------------------------------------------

func TestUserService_ChangePassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newService(t)

		existing := validUser()
		existing.Password = "OldHash"

		mockRepo.EXPECT().GetUser(testCtx, existing.Id).Return(existing, nil)
		// After ChangePassword the user object passed to repo has a bcrypt hash of the new password.
		mockRepo.EXPECT().
			UpdateUser(testCtx, mock.MatchedBy(func(u domain.User) bool {
				return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("NewPass1")) == nil
			}), "admin").
			Return(existing, nil)

		_, err := svc.ChangePassword(testCtx, existing.Id, "NewPass1", "admin")
		require.NoError(t, err)
	})

	t.Run("user_not_found", func(t *testing.T) {
		svc, mockRepo := newService(t)

		mockRepo.EXPECT().GetUser(testCtx, 99).Return(domain.User{}, errors.New("not found"))

		_, err := svc.ChangePassword(testCtx, 99, "NewPass1", "admin")
		assert.Error(t, err)
	})

	t.Run("weak_password_skips_repo_update", func(t *testing.T) {
		svc, mockRepo := newService(t)

		existing := validUser()
		mockRepo.EXPECT().GetUser(testCtx, existing.Id).Return(existing, nil)

		_, err := svc.ChangePassword(testCtx, existing.Id, "weak", "admin")
		assert.Error(t, err)
	})

	t.Run("repo_update_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newService(t)

		existing := validUser()
		mockRepo.EXPECT().GetUser(testCtx, existing.Id).Return(existing, nil)
		mockRepo.EXPECT().UpdateUser(testCtx, mock.Anything, "admin").Return(domain.User{}, errors.New("db error"))

		_, err := svc.ChangePassword(testCtx, existing.Id, "NewPass1", "admin")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// DeleteUser
// ---------------------------------------------------------------------------

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newService(t)

		mockRepo.EXPECT().DeleteUser(testCtx, 1, "admin").Return(nil)

		err := svc.DeleteUser(testCtx, 1, "admin")
		assert.NoError(t, err)
	})

	t.Run("repo_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newService(t)

		mockRepo.EXPECT().DeleteUser(testCtx, 1, "admin").Return(errors.New("db error"))

		err := svc.DeleteUser(testCtx, 1, "admin")
		assert.Error(t, err)
	})
}
