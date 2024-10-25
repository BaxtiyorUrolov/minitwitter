package service_test

import (
	"context"
	"errors"
	"testing"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/pkg/security"
	"twitter/service"
	"twitter/storage/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockUserStorage := new(mocks.MockUserStorage)
	mockLogger := logger.New("test_service")

	// Misol uchun mavjud foydalanuvchi
	hashedPassword, err := security.HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user := models.User{
		ID:       "user-123",
		UserName: "testuser",
		Password: hashedPassword, // hashed parolni olingan qiymatga bog'lash
	}

	mockStorage.On("User").Return(mockUserStorage)
	mockUserStorage.On("GetUserCredentials", mock.Anything, user.UserName).Return(user, nil)

	authService := service.NewAuthService(mockStorage, mockLogger)

	t.Run("Success", func(t *testing.T) {
		loginRequest := models.LoginRequest{
			Login:    user.UserName,
			Password: "password123",
		}

		ctx := context.Background()
		resp, err := authService.Login(ctx, loginRequest)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		loginRequest := models.LoginRequest{
			Login:    user.UserName,
			Password: "wrongpassword",
		}

		ctx := context.Background()
		resp, err := authService.Login(ctx, loginRequest)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("User Not Found", func(t *testing.T) {
		loginRequest := models.LoginRequest{
			Login:    "nonexistentuser",
			Password: "password123",
		}

		mockUserStorage.On("GetUserCredentials", mock.Anything, "nonexistentuser").Return(models.User{}, errors.New("user not found"))

		ctx := context.Background()
		resp, err := authService.Login(ctx, loginRequest)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})
}
