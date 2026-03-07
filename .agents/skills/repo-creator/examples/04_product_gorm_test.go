package repo

import (
	"testing"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/pkg/ctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductGormRepo_Create(t *testing.T) {
	db, err := NewGormRepo(TEST_CONN_STRING, 1)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		repo := NewProductGormRepo(tx)

		product := domain.Product{
			Name:        "Test Product",
			Description: "A product for testing",
			Price:       99.99,
		}

		createdProduct, err := repo.InsertProduct(ctx.NewCtx("111111"), product, "sysadmin")
		require.NoError(t, err)
		assert.NotZero(t, createdProduct.Id, "Id should be set upon creation")
		assert.Equal(t, "Test Product", createdProduct.Name)
		assert.Equal(t, 99.99, createdProduct.Price)
	})
}
