package service

import (
	"context"
	"testing"
	"twitter/api/models"
	"twitter/config"
	"twitter/pkg/logger"
	"twitter/storage/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Create(t *testing.T) {

	mockStorage := new(mocks.MockStorage)
	mockUserStorage := new(mocks.MockUserStorage)
	mockLogger := logger.New(config.Load().ServiceName)

	mockStorage.On("User").Return(mockUserStorage)
	mockUserStorage.On("Create", mock.Anything, mock.Anything).Return("1", nil)
	mockUserStorage.On("GetByID", mock.Anything, models.PrimaryKey{ID: "1"}).Return(models.User{
		ID:       "1",
		UserName: "testuser",
	}, nil)

	userService := NewUserService(mockStorage, mockLogger)
	createUser := models.CreateUser{UserName: "testuser", Password: "password123"}

	user, err := userService.Create(context.Background(), createUser)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.UserName)
	mockStorage.AssertExpectations(t)
	mockUserStorage.AssertExpectations(t)
}

func TestUserService_IsLoginExist(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockUserStorage := new(mocks.MockUserStorage)
	mockLogger := logger.New(config.Load().ServiceName)

	mockStorage.On("User").Return(mockUserStorage)
	mockUserStorage.On("IsUserNameExist", mock.Anything, "testuser").Return(true, nil)

	userService := NewUserService(mockStorage, mockLogger)
	login := "testuser"

	exists, err := userService.IsLoginExist(context.Background(), login)

	assert.NoError(t, err)
	assert.True(t, exists)
	mockStorage.AssertExpectations(t)
	mockUserStorage.AssertExpectations(t)
}

func TestUserService_GetByID(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockUserStorage := new(mocks.MockUserStorage)
	mockLogger := logger.New(config.Load().ServiceName)

	mockStorage.On("User").Return(mockUserStorage)
	mockUserStorage.On("GetByID", mock.Anything, models.PrimaryKey{ID: "1"}).Return(models.User{
		ID:       "1",
		UserName: "testuser",
	}, nil)

	userService := NewUserService(mockStorage, mockLogger)
	userID := models.PrimaryKey{ID: "1"}

	user, err := userService.GetByID(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.UserName)
	mockStorage.AssertExpectations(t)
	mockUserStorage.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockUserStorage := new(mocks.MockUserStorage)
	mockLogger := logger.New(config.Load().ServiceName)

	mockStorage.On("User").Return(mockUserStorage)
	mockUserStorage.On("Update", mock.Anything, mock.Anything).Return(nil)

	userService := NewUserService(mockStorage, mockLogger)
	updateUser := models.UpdateUser{
		ID:       "1",
		UserName: "updatedUser",
	}

	err := userService.Update(context.Background(), updateUser)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
	mockUserStorage.AssertExpectations(t)
}

func TestUserService_Delete(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockUserStorage := new(mocks.MockUserStorage)
	mockLogger := logger.New(config.Load().ServiceName)

	mockStorage.On("User").Return(mockUserStorage)
	mockUserStorage.On("Delete", mock.Anything, mock.Anything).Return(nil)

	userService := NewUserService(mockStorage, mockLogger)
	userID := models.PrimaryKey{ID: "1"}

	err := userService.Delete(context.Background(), userID)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
	mockUserStorage.AssertExpectations(t)
}

func TestUserService_GetAllUsers(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockUserStorage := new(mocks.MockUserStorage)
	mockLogger := logger.New(config.Load().ServiceName)

	mockStorage.On("User").Return(mockUserStorage)
	mockUserStorage.On("GetList", mock.Anything, mock.Anything).Return(models.UsersResponse{
		Users: []models.User{
			{ID: "1", UserName: "user1"},
			{ID: "2", UserName: "user2"},
		},
		Count: 2,
	}, nil)

	userService := NewUserService(mockStorage, mockLogger)
	request := models.GetListRequest{Page: 1, Limit: 10}

	usersResponse, err := userService.GetAllUsers(context.Background(), request)

	assert.NoError(t, err)
	assert.Len(t, usersResponse.Users, 2)
	mockStorage.AssertExpectations(t)
	mockUserStorage.AssertExpectations(t)
}
