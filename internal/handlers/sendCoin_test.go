package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go/avito-test/pkg/midleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCoinService is a mock implementation of CoinService
type MockCoinService struct {
	mock.Mock
}

func (m *MockCoinService) SendCoins(senderName string, receiverUsername string, amount int) error {
	args := m.Called(senderName, receiverUsername, amount)
	return args.Error(0)
}

func TestSendCoinHandler_SendCoin(t *testing.T) {
	// Create mock dependencies
	mockCoinService := new(MockCoinService)

	// Create the handler with mock dependencies
	handler := &SendCoinHandler{
		CoinService: mockCoinService,
	}

	// Test case 1: Successful coin transfer
	t.Run("Success", func(t *testing.T) {
		// Mock dependencies
		mockCoinService.On("SendCoins", "sender", "receiver", 100).Return(nil)

		// Create a request with a valid context and body
		reqBody := `{"toUser": "receiver", "amount": 100}`
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), midleware.ContextUserNameKey, "sender")
		req = req.WithContext(ctx)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the handler
		handler.SendCoin()(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Coins sent successfully"}`, w.Body.String())
	})

	// Test case 2: Invalid request body
	t.Run("Invalid Request Body", func(t *testing.T) {
		// Create a request with an invalid body
		reqBody := `{"toUser": "receiver", "amount": "invalid"}` // Invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), midleware.ContextUserNameKey, "sender")
		req = req.WithContext(ctx)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the handler
		handler.SendCoin()(w, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "of type int")
	})

	// Test case: Coin transfer failure
	t.Run("Coin Transfer Failure", func(t *testing.T) {
		// Mock dependencies
		mockCoinService.On("SendCoins", "sender", "receiver", 100).Return(errors.New("insufficient coins"))

		// Create a request with a valid context and body
		reqBody := `{"toUser": "receiver", "amount": 100}`
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), midleware.ContextUserNameKey, "sender")
		req = req.WithContext(ctx)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the handler
		handler.SendCoin()(w, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "insufficient coins")
	})
}
