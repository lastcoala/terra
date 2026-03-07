package domain

import (
	"fmt"
	"strings"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// ValidateName returns an error if the product name is empty.
func (p Product) ValidateName() error {
	if strings.TrimSpace(p.Name) == "" {
		return p.error(fmt.Errorf("name must not be empty"), "ValidateName")
	}
	return nil
}

// ValidatePrice returns an error if the price is not positive.
func (p Product) ValidatePrice() error {
	if p.Price <= 0 {
		return p.error(fmt.Errorf("price must be greater than zero, got %v", p.Price), "ValidatePrice")
	}
	return nil
}

// ValidateStock returns an error if the stock quantity is negative.
func (p Product) ValidateStock() error {
	if p.Stock < 0 {
		return p.error(fmt.Errorf("stock must be non-negative, got %v", p.Stock), "ValidateStock")
	}
	return nil
}

// Validate returns an error if any field of the product is invalid.
func (p Product) Validate() error {
	if err := p.ValidateName(); err != nil {
		return p.error(err, "Validate")
	}
	if err := p.ValidatePrice(); err != nil {
		return p.error(err, "Validate")
	}
	if err := p.ValidateStock(); err != nil {
		return p.error(err, "Validate")
	}
	return nil
}

func (p Product) error(err error, method string, params ...any) error {
	return fmt.Errorf("Product.(%v)(%v) %w", method, params, err)
}
