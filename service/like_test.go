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

func setupLikeService() (likeService, *mocks.MockStorage, *logger.Logger) {
	mockStorage := new(mocks.MockStorage)
	mockLogger := logger.New("test_service")
	return NewLikeService(mockStorage, mockLogger), mockStorage, &mockLogger
}

func TestLikeService_LikeTweet(t *testing.T) {
	likeService, mockStorage, _ := setupLikeService()
	like := models.Like{UserID: "user1", TweetID: "tweet1"}

	mockStorage.On("Like").Return(mockStorage)
	mockStorage.On("LikeTweet", mock.Anything, like).Return(nil)

	err := likeService.LikeTweet(context.Background(), like)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestLikeService_UnlikeTweet(t *testing.T) {
	likeService, mockStorage, _ := setupLikeService()
	userID := "user1"
	tweetID := "tweet1"

	mockStorage.On("Like").Return(mockStorage)
	mockStorage.On("UnlikeTweet", mock.Anything, userID, tweetID).Return(nil)

	err := likeService.UnlikeTweet(context.Background(), userID, tweetID)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestLikeService_GetLikeCount(t *testing.T) {
	likeService, mockStorage, _ := setupLikeService()
	tweetID := "tweet1"
	expectedCount := 10

	mockStorage.On("Like").Return(mockStorage)
	mockStorage.On("GetLikeCount", mock.Anything, tweetID).Return(expectedCount, nil)

	count, err := likeService.GetLikeCount(context.Background(), tweetID)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockStorage.AssertExpectations(t)
}
