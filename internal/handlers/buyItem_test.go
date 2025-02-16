package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go/avito-test/pkg/midleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderService is a mock implementation of OrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) PurchaseItem(username string, itemName string) error {
	args := m.Called(username, itemName)
	return args.Error(0)
}

func TestBuyItemHandler_BuyItem(t *testing.T) {
	// Create mock dependencies
	mockOrderService := new(MockOrderService)

	// Create the handler with mock dependencies
	handler := &BuyItemHandler{
		OrderService: mockOrderService,
	}

	// Test case 1: Successful item purchase
	t.Run("Success", func(t *testing.T) {
		// Mock dependencies
		mockOrderService.On("PurchaseItem", "buyer", "item1").Return(nil)

		// Create a request with a valid context and path parameter
		req := httptest.NewRequest(http.MethodGet, "/api/buy/item1", nil)
		req.SetPathValue("item", "item1")
		ctx := context.WithValue(req.Context(), midleware.ContextUserNameKey, "buyer")
		req = req.WithContext(ctx)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the handler
		handler.BuyItem()(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Item purchased successfully"}`, w.Body.String())
	})

	// Test case 2: Item purchase failure
	t.Run("Purchase Failure", func(t *testing.T) {
		// Mock dependencies
		mockOrderService.On("PurchaseItem", "buyer", "item1").Return(errors.New("item out of stock"))

		// Create a request with a valid context and path parameter
		req := httptest.NewRequest(http.MethodGet, "/api/buy/item1", nil)
		req.SetPathValue("item", "item1")
		ctx := context.WithValue(req.Context(), midleware.ContextUserNameKey, "buyer")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.BuyItem()(w, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "item out of stock")
	})
}
