package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Product struct {
	ID   int
	Name string
}

type ProductRepository interface {
	GetProductByID(id int) (*Product, error)
	GetAllProducts() ([]*Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func (s *ProductService) GetProductByID(id int) (*Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *ProductService) GetAllProducts() ([]*Product, error) {
	return s.repo.GetAllProducts()
}

type MockProductRepository struct {
	products map[int]*Product
	err      error
}

func (m *MockProductRepository) GetProductByID(id int) (*Product, error) {
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, m.err
}

func (m *MockProductRepository) GetAllProducts() ([]*Product, error) {
	if m.products == nil {
		return nil, m.err
	}
	products := make([]*Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func TestProductService_GetProductByID(t *testing.T) {
	mockRepo := &MockProductRepository{
		products: map[int]*Product{
			1: {ID: 1, Name: "Test Product"},
		},
		err: errors.New("Product not found"),
	}
	service := &ProductService{repo: mockRepo}

	t.Run("Product Found", func(t *testing.T) {
		expected := &Product{ID: 1, Name: "Test Product"}
		actual, err := service.GetProductByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Product Not Found", func(t *testing.T) {
		expectedErr := errors.New("Product not found")
		actual, err := service.GetProductByID(2)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, expectedErr, err)
	})
}

func TestProductService_GetAllProducts(t *testing.T) {
	mockRepo := &MockProductRepository{
		products: map[int]*Product{
			1: {ID: 1, Name: "Product 1"},
			2: {ID: 2, Name: "Product 2"},
		},
		err: errors.New("No products found"),
	}
	service := &ProductService{repo: mockRepo}

	t.Run("Products Found", func(t *testing.T) {
		expected := []*Product{
			{ID: 1, Name: "Product 1"},
			{ID: 2, Name: "Product 2"},
		}
		actual, err := service.GetAllProducts()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Products Not Found", func(t *testing.T) {
		mockRepo := &MockProductRepository{err: errors.New("No products found")}
		service := &ProductService{repo: mockRepo}
		expectedErr := errors.New("No products found")
		actual, err := service.GetAllProducts()
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, expectedErr, err)
	})
}
