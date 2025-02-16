package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go/avito-test/internal/models"
	"go/avito-test/pkg/midleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserInfo(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserService) FindOrCreateUser(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserService) GetUserByName(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestInfoHandler_Info(t *testing.T) {
	// Create mock dependencies
	mockUserService := new(MockUserService)

	// Create the handler with mock dependencies
	handler := &InfoHandler{
		UserService: mockUserService,
	}

	// Test case 1: Successful retrieval of user info
	t.Run("Success", func(t *testing.T) {
		// Mock dependencies
		mockUserService.On("GetUserInfo", "testuser").Return(&models.User{
			Username:                 "testuser",
			CoinBalance:              1000,
			Orders:                   []models.Order{},
			ReceivedCoinTransactions: []models.CoinTransaction{},
			SentCoinTransactions:     []models.CoinTransaction{},
		}, nil)

		// Create a request with a valid context
		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		ctx := context.WithValue(req.Context(), midleware.ContextUserNameKey, "testuser")
		req = req.WithContext(ctx)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the handler
		handler.Info()(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
            "coins": 1000,
            "inventory": [],
            "coinHistory": {
                "received": [],
                "sent": []
            }
        }`, w.Body.String())
	})

	// Test case 2: Missing username in context
	t.Run("Missing Username", func(t *testing.T) {
		// Create a request without a username in the context
		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the handler
		handler.Info()(w, req)

		// Assert the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "username not found in context")
	})

	// Test case 3: User info retrieval failure
	t.Run("User Info Retrieval Failure", func(t *testing.T) {
		// Mock dependencies
		mockUserService.On("GetUserInfo", "testuser").Return(&models.User{}, errors.New("database error"))

		// Create a request with a valid context
		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		ctx := context.WithValue(req.Context(), midleware.ContextUserNameKey, "testuser")
		req = req.WithContext(ctx)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the handler
		handler.Info()(w, req)

		// Assert the response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}
