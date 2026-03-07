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
)

var testCtx = ctx.NewCtx("1234567890")

// newProductService wires a fresh MockIProductRepo into a ProductService.
func newProductService(t *testing.T) (*ProductService, *mocks.MockIProductRepo) {
	t.Helper()
	mockRepo := mocks.NewMockIProductRepo(t)
	svc := NewProductService(mockRepo)
	return svc, mockRepo
}

// validProduct returns a domain.Product that satisfies all validation rules.
func validProduct() domain.Product {
	return domain.Product{
		Id:          1,
		Name:        "Mechanical Keyboard",
		Description: "A great keyboard",
		Price:       79.99,
		Stock:       50,
	}
}

// ---------------------------------------------------------------------------
// InsertProduct
// ---------------------------------------------------------------------------

func TestProductService_InsertProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		product := validProduct()
		product.Id = 0
		expected := validProduct()

		mockRepo.EXPECT().InsertProduct(testCtx, product, "sysadmin").Return(expected, nil)

		created, err := svc.InsertProduct(testCtx, product, "sysadmin")
		require.NoError(t, err)
		assert.Equal(t, expected, created)
	})

	t.Run("invalid_name_skips_repo", func(t *testing.T) {
		svc, _ := newProductService(t)

		product := validProduct()
		product.Name = ""

		_, err := svc.InsertProduct(testCtx, product, "sysadmin")
		assert.Error(t, err)
	})

	t.Run("invalid_price_skips_repo", func(t *testing.T) {
		svc, _ := newProductService(t)

		product := validProduct()
		product.Price = 0

		_, err := svc.InsertProduct(testCtx, product, "sysadmin")
		assert.Error(t, err)
	})

	t.Run("repo_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		product := validProduct()
		mockRepo.EXPECT().InsertProduct(testCtx, mock.Anything, "sysadmin").Return(domain.Product{}, errors.New("db error"))

		_, err := svc.InsertProduct(testCtx, product, "sysadmin")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// GetProduct
// ---------------------------------------------------------------------------

func TestProductService_GetProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		expected := validProduct()
		mockRepo.EXPECT().GetProduct(testCtx, 1).Return(expected, nil)

		got, err := svc.GetProduct(testCtx, 1)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not_found_propagates_error", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		mockRepo.EXPECT().GetProduct(testCtx, 999).Return(domain.Product{}, errors.New("record not found"))

		_, err := svc.GetProduct(testCtx, 999)
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// GetProducts
// ---------------------------------------------------------------------------

func TestProductService_GetProducts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		expected := []domain.Product{validProduct()}
		mockRepo.EXPECT().GetProducts(testCtx, 0, 10).Return(expected, nil)

		got, err := svc.GetProducts(testCtx, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("empty_result", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		mockRepo.EXPECT().GetProducts(testCtx, 0, 10).Return([]domain.Product{}, nil)

		got, err := svc.GetProducts(testCtx, 0, 10)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("repo_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		mockRepo.EXPECT().GetProducts(testCtx, 0, 10).Return(nil, errors.New("db error"))

		_, err := svc.GetProducts(testCtx, 0, 10)
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// UpdateProduct
// ---------------------------------------------------------------------------

func TestProductService_UpdateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		existing := validProduct()
		input := existing
		input.Name = "Updated Keyboard"
		expected := input

		mockRepo.EXPECT().GetProduct(testCtx, existing.Id).Return(existing, nil)
		mockRepo.EXPECT().UpdateProduct(testCtx, mock.MatchedBy(func(p domain.Product) bool {
			return p.Id == existing.Id && p.Name == "Updated Keyboard"
		}), "admin").Return(expected, nil)

		got, err := svc.UpdateProduct(testCtx, input, "admin")
		require.NoError(t, err)
		assert.Equal(t, "Updated Keyboard", got.Name)
	})

	t.Run("product_not_found", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		mockRepo.EXPECT().GetProduct(testCtx, 99).Return(domain.Product{}, errors.New("not found"))

		product := validProduct()
		product.Id = 99
		_, err := svc.UpdateProduct(testCtx, product, "admin")
		assert.Error(t, err)
	})

	t.Run("invalid_product_skips_repo_update", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		existing := validProduct()
		mockRepo.EXPECT().GetProduct(testCtx, existing.Id).Return(existing, nil)

		bad := existing
		bad.Price = -10
		_, err := svc.UpdateProduct(testCtx, bad, "admin")
		assert.Error(t, err)
	})

	t.Run("repo_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		existing := validProduct()
		mockRepo.EXPECT().GetProduct(testCtx, existing.Id).Return(existing, nil)
		mockRepo.EXPECT().UpdateProduct(testCtx, mock.Anything, "admin").Return(domain.Product{}, errors.New("db error"))

		_, err := svc.UpdateProduct(testCtx, existing, "admin")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// DeleteProduct
// ---------------------------------------------------------------------------

func TestProductService_DeleteProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		mockRepo.EXPECT().DeleteProduct(testCtx, 1, "admin").Return(nil)

		err := svc.DeleteProduct(testCtx, 1, "admin")
		assert.NoError(t, err)
	})

	t.Run("repo_error_is_propagated", func(t *testing.T) {
		svc, mockRepo := newProductService(t)

		mockRepo.EXPECT().DeleteProduct(testCtx, 1, "admin").Return(errors.New("db error"))

		err := svc.DeleteProduct(testCtx, 1, "admin")
		assert.Error(t, err)
	})
}
