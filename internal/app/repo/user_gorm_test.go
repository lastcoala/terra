package repo

import (
	"testing"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/pkg/ctx"
	"github.com/lastcoala/terra/pkg/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// seedUser is a helper that inserts one user inside the supplied repo and
// returns the created domain.User (with an auto-assigned Id).
func seedUser(t *testing.T, repo IUserRepo, u domain.User) domain.User {
	t.Helper()
	created, err := repo.InsertUser(ctx.NewCtx("seed"), u, "sysadmin")
	require.NoError(t, err)
	return created
}

// ---------------------------------------------------------------------------
// InsertUser
// ---------------------------------------------------------------------------

func TestUserGormRepo_InsertUser(t *testing.T) {
	db, err := NewGormRepo(TEST_CONN_STRING, 1)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		user := domain.User{
			Name:     "Adrian",
			Gender:   "male",
			Email:    "adrianch@example.com",
			Password: "hashedpassword",
		}

		created, err := repo.InsertUser(ctx.NewCtx("111111"), user, "sysadmin")
		require.NoError(t, err)
		assert.NotZero(t, created.Id, "Id should be auto-set")
		assert.Equal(t, user.Name, created.Name)
		assert.Equal(t, user.Gender, created.Gender)
		assert.Equal(t, user.Email, created.Email)
		assert.Equal(t, user.Password, created.Password)
	})

	t.Run("duplicate_email_fails", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		user := domain.User{
			Name:     "Adrian",
			Gender:   "male",
			Email:    "duplicate@example.com",
			Password: "hashedpassword",
		}

		_, err := repo.InsertUser(ctx.NewCtx("222222"), user, "sysadmin")
		require.NoError(t, err)

		_, err = repo.InsertUser(ctx.NewCtx("333333"), user, "sysadmin")
		require.Error(t, err, "should fail on duplicate email due to unique index")
	})
}

// ---------------------------------------------------------------------------
// GetUser
// ---------------------------------------------------------------------------

func TestUserGormRepo_GetUser(t *testing.T) {
	db, err := NewGormRepo(TEST_CONN_STRING, 1)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		seeded := seedUser(t, repo, domain.User{
			Name:     "Alice",
			Gender:   "female",
			Email:    "alice@example.com",
			Password: "pass",
		})

		got, err := repo.GetUser(ctx.NewCtx("111"), seeded.Id)
		require.NoError(t, err)
		assert.Equal(t, seeded.Id, got.Id)
		assert.Equal(t, "Alice", got.Name)
		assert.Equal(t, "female", got.Gender)
		assert.Equal(t, "alice@example.com", got.Email)
	})

	t.Run("not_found", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		_, err := repo.GetUser(ctx.NewCtx("999"), 999999)
		require.Error(t, err, "should return error for non-existent id")
	})
}

// ---------------------------------------------------------------------------
// GetUsers
// ---------------------------------------------------------------------------

func TestUserGormRepo_GetUsers(t *testing.T) {
	db, err := NewGormRepo(TEST_CONN_STRING, 1)
	require.NoError(t, err)

	t.Run("returns_all_without_filters", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		seedUser(t, repo, domain.User{Name: "Bob", Gender: "male", Email: "bob@example.com", Password: "p1"})
		seedUser(t, repo, domain.User{Name: "Carol", Gender: "female", Email: "carol@example.com", Password: "p2"})

		users, err := repo.GetUsers(ctx.NewCtx("aaa"), 0, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 2)
	})

	t.Run("filter_by_gender", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		seedUser(t, repo, domain.User{Name: "Dave", Gender: "male", Email: "dave@example.com", Password: "p3"})
		seedUser(t, repo, domain.User{Name: "Eve", Gender: "female", Email: "eve@example.com", Password: "p4"})

		users, err := repo.GetUsers(ctx.NewCtx("bbb"), 0, 10, filter.Filter{
			Attribute: "gender",
			Operator:  "=",
			Value:     "male",
		})
		require.NoError(t, err)
		for _, u := range users {
			assert.Equal(t, "male", u.Gender)
		}
	})

	t.Run("offset_and_limit", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		seedUser(t, repo, domain.User{Name: "F1", Gender: "male", Email: "f1@example.com", Password: "p"})
		seedUser(t, repo, domain.User{Name: "F2", Gender: "male", Email: "f2@example.com", Password: "p"})
		seedUser(t, repo, domain.User{Name: "F3", Gender: "male", Email: "f3@example.com", Password: "p"})

		all, err := repo.GetUsers(ctx.NewCtx("ccc"), 0, 100)
		require.NoError(t, err)

		paged, err := repo.GetUsers(ctx.NewCtx("ddd"), 0, 1)
		require.NoError(t, err)
		assert.Equal(t, 1, len(paged))
		assert.Less(t, len(paged), len(all))
	})

	t.Run("empty_result", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		users, err := repo.GetUsers(ctx.NewCtx("eee"), 0, 10, filter.Filter{
			Attribute: "email",
			Operator:  "=",
			Value:     "nonexistent@nowhere.com",
		})
		require.NoError(t, err)
		assert.Empty(t, users)
	})
}

// ---------------------------------------------------------------------------
// UpdateUser
// ---------------------------------------------------------------------------

func TestUserGormRepo_UpdateUser(t *testing.T) {
	db, err := NewGormRepo(TEST_CONN_STRING, 1)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		seeded := seedUser(t, repo, domain.User{
			Name:     "OldName",
			Gender:   "male",
			Email:    "update@example.com",
			Password: "oldpass",
		})

		// Capture the CreatedAt before the update so we can compare timestamps.
		var before UserGormModel
		require.NoError(t, tx.First(&before, seeded.Id).Error)

		seeded.Name = "NewName"
		seeded.Password = "newpass"

		updated, err := repo.UpdateUser(ctx.NewCtx("upd"), seeded, "admin")
		require.NoError(t, err)
		assert.Equal(t, seeded.Id, updated.Id)
		assert.Equal(t, "NewName", updated.Name)
		assert.Equal(t, "newpass", updated.Password)

		// Verify domain fields are persisted within the transaction.
		got, err := repo.GetUser(ctx.NewCtx("upd2"), seeded.Id)
		require.NoError(t, err)
		assert.Equal(t, "NewName", got.Name)

		// Verify BaseModel audit fields via a direct model query on the same tx.
		var model UserGormModel
		require.NoError(t, tx.First(&model, seeded.Id).Error)
		assert.Equal(t, "admin", model.UpdatedBy, "UpdatedBy should be set to the updater")
		assert.True(t, model.UpdatedAt.After(before.CreatedAt) || model.UpdatedAt.Equal(before.CreatedAt),
			"UpdatedAt (%v) should be >= CreatedAt (%v)", model.UpdatedAt, before.CreatedAt)
	})

	t.Run("empty_password_does_not_overwrite", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		seeded := seedUser(t, repo, domain.User{
			Name:     "Bob",
			Gender:   "male",
			Email:    "bob@example.com",
			Password: "originalpass",
		})

		// Update name only; leave Password empty.
		seeded.Name = "BobUpdated"
		seeded.Password = ""

		updated, err := repo.UpdateUser(ctx.NewCtx("upd3"), seeded, "admin")
		require.NoError(t, err)
		assert.Equal(t, "BobUpdated", updated.Name)

		// The stored password should still be the original value.
		var model UserGormModel
		require.NoError(t, tx.First(&model, seeded.Id).Error)
		assert.Equal(t, "originalpass", model.Password, "password should not be overwritten when empty")
	})
}

// ---------------------------------------------------------------------------
// DeleteUser (soft delete)
// ---------------------------------------------------------------------------

func TestUserGormRepo_DeleteUser(t *testing.T) {
	db, err := NewGormRepo(TEST_CONN_STRING, 1)
	require.NoError(t, err)

	t.Run("soft_delete_sets_is_deleted", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Use raw gorm for verification (soft delete is a model-level flag, not visible via repo)
		repo := NewUserGormRepo(tx)

		seeded := seedUser(t, repo, domain.User{
			Name:     "ToDelete",
			Gender:   "male",
			Email:    "todelete@example.com",
			Password: "pass",
		})

		err := repo.DeleteUser(ctx.NewCtx("del"), seeded.Id, "admin")
		require.NoError(t, err)

		// Verify is_deleted flag via direct model query on the same tx
		var model UserGormModel
		err = tx.Unscoped().First(&model, seeded.Id).Error
		require.NoError(t, err)
		assert.True(t, model.IsDeleted, "is_deleted should be true after soft delete")
		assert.Equal(t, "admin", model.UpdatedBy)
	})

	t.Run("delete_non_existent_id_returns_no_error", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewUserGormRepo(tx)

		// GORM Updates on zero rows is not an error by default
		err := repo.DeleteUser(ctx.NewCtx("del2"), 999999, "admin")
		assert.NoError(t, err)
	})
}
