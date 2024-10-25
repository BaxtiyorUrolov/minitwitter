package service

import (
	"context"
	"testing"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/storage/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupFollowService() (followService, *mocks.MockStorage, *logger.Logger) {
	mockStorage := new(mocks.MockStorage)
	mockLogger := logger.New("test_service")
	return NewFollowService(mockStorage, mockLogger), mockStorage, &mockLogger
}

func TestFollowService_FollowUser(t *testing.T) {
	followService, mockStorage, _ := setupFollowService()
	follow := models.Follow{FollowerID: "user1", FollowingID: "user2"}

	mockStorage.On("Follow").Return(mockStorage)
	mockStorage.On("FollowUser", mock.Anything, follow).Return(nil)

	err := followService.FollowUser(context.Background(), follow)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestFollowService_UnfollowUser(t *testing.T) {
	followService, mockStorage, _ := setupFollowService()
	followerID := "user1"
	followingID := "user2"

	mockStorage.On("Follow").Return(mockStorage)
	mockStorage.On("UnfollowUser", mock.Anything, followerID, followingID).Return(nil)

	err := followService.UnfollowUser(context.Background(), followerID, followingID)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestFollowService_GetFollowers(t *testing.T) {
	followService, mockStorage, _ := setupFollowService()
	userID := "user1"
	expectedFollowers := []models.Follow{
		{FollowerID: "user2", FollowingID: "user1"},
	}

	mockStorage.On("Follow").Return(mockStorage)
	mockStorage.On("GetFollowers", mock.Anything, userID).Return(expectedFollowers, nil)

	followers, err := followService.GetFollowers(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedFollowers, followers)
	mockStorage.AssertExpectations(t)
}

func TestFollowService_GetFollowings(t *testing.T) {
	followService, mockStorage, _ := setupFollowService()
	userID := "user1"
	expectedFollowings := []models.Follow{
		{FollowerID: "user1", FollowingID: "user2"},
	}

	mockStorage.On("Follow").Return(mockStorage)
	mockStorage.On("GetFollowing", mock.Anything, userID).Return(expectedFollowings, nil)

	followings, err := followService.GetFollowings(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedFollowings, followings)
	mockStorage.AssertExpectations(t)
}
