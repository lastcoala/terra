package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// ValidateName
// ---------------------------------------------------------------------------

func TestProduct_ValidateName(t *testing.T) {
	tests := []struct {
		name    string
		pName   string
		wantErr bool
	}{
		{"valid name", "Keyboard", false},
		{"empty string", "", true},
		{"whitespace only", "   ", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := Product{Name: tc.pName}
			err := p.ValidateName()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ValidatePrice
// ---------------------------------------------------------------------------

func TestProduct_ValidatePrice(t *testing.T) {
	tests := []struct {
		name    string
		price   float64
		wantErr bool
	}{
		{"valid price", 9.99, false},
		{"zero price", 0, true},
		{"negative price", -1.0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := Product{Price: tc.price}
			err := p.ValidatePrice()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ValidateStock
// ---------------------------------------------------------------------------

func TestProduct_ValidateStock(t *testing.T) {
	tests := []struct {
		name    string
		stock   int
		wantErr bool
	}{
		{"valid zero stock", 0, false},
		{"valid positive stock", 100, false},
		{"negative stock", -1, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := Product{Stock: tc.stock}
			err := p.ValidateStock()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Validate
// ---------------------------------------------------------------------------

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		wantErr bool
	}{
		{
			name:    "valid product",
			product: Product{Name: "Keyboard", Price: 49.99, Stock: 10},
			wantErr: false,
		},
		{
			name:    "empty name",
			product: Product{Name: "", Price: 49.99, Stock: 10},
			wantErr: true,
		},
		{
			name:    "zero price",
			product: Product{Name: "Keyboard", Price: 0, Stock: 10},
			wantErr: true,
		},
		{
			name:    "negative stock",
			product: Product{Name: "Keyboard", Price: 49.99, Stock: -1},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.product.Validate()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
