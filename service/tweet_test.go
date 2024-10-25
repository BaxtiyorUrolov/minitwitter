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

func setupTweetService() (tweetService, *mocks.MockStorage, *logger.Logger) {
	mockStorage := new(mocks.MockStorage)
	mockLogger := logger.New(config.Load().ServiceName)

	return NewTweetService(mockStorage, mockLogger), mockStorage, &mockLogger
}

func TestTweetService_Create(t *testing.T) {

	tweetService, mockStorage, _ := setupTweetService()
	createTweet := models.CreateTweet{Content: "Test tweet"}

	mockStorage.On("Tweet").Return(mockStorage)
	mockStorage.On("CreateTweet", mock.Anything, createTweet).Return("1", nil)

	err := tweetService.Create(context.Background(), createTweet)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestTweetService_Get(t *testing.T) {
	tweetService, mockStorage, _ := setupTweetService()
	tweetID := "1"
	expectedTweet := models.Tweet{ID: tweetID, Content: "Test tweet"}

	mockStorage.On("Tweet").Return(mockStorage)
	mockStorage.On("GetTweet", mock.Anything, tweetID).Return(expectedTweet, nil)

	tweet, err := tweetService.Get(context.Background(), tweetID)
	assert.NoError(t, err)
	assert.Equal(t, expectedTweet, tweet)
	mockStorage.AssertExpectations(t)
}

func TestTweetService_Update(t *testing.T) {
	tweetService, mockStorage, _ := setupTweetService()
	updateTweet := models.UpdateTweet{ID: "1", Content: "Updated tweet"}

	mockStorage.On("Tweet").Return(mockStorage)
	mockStorage.On("UpdateTweet", mock.Anything, updateTweet).Return(nil)

	err := tweetService.Update(context.Background(), updateTweet)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestTweetService_Delete(t *testing.T) {
	tweetService, mockStorage, _ := setupTweetService()
	tweetID := "1"

	mockStorage.On("Tweet").Return(mockStorage)
	mockStorage.On("DeleteTweet", mock.Anything, tweetID).Return(nil)

	err := tweetService.Delete(context.Background(), tweetID)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestTweetService_GetTweetsByUser(t *testing.T) {
	tweetService, mockStorage, _ := setupTweetService()
	userID := "1"
	expectedResponse := models.TweetsResponse{Count: 1, Tweets: []models.Tweet{{ID: "1", Content: "User's tweet"}}}

	mockStorage.On("Tweet").Return(mockStorage)
	mockStorage.On("ListTweetsByUser", mock.Anything, userID).Return(expectedResponse, nil)

	tweets, err := tweetService.GetTweetsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, tweets)
	mockStorage.AssertExpectations(t)
}

func TestTweetService_GetList(t *testing.T) {
	tweetService, mockStorage, _ := setupTweetService()
	request := models.GetListRequest{Page: 1, Limit: 10}
	expectedResponse := models.TweetsResponse{Count: 1, Tweets: []models.Tweet{{ID: "1", Content: "Listed tweet"}}}

	mockStorage.On("Tweet").Return(mockStorage)
	mockStorage.On("GetTweetList", mock.Anything, request).Return(expectedResponse, nil)

	tweets, err := tweetService.GetList(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, tweets)
	mockStorage.AssertExpectations(t)
}

func TestTweetService_IsTweetOwner(t *testing.T) {
	tweetService, mockStorage, _ := setupTweetService()
	tweetID := "1"
	userID := "user1"
	tweet := models.Tweet{ID: tweetID, UserID: userID}

	mockStorage.On("Tweet").Return(mockStorage)
	mockStorage.On("GetTweet", mock.Anything, tweetID).Return(tweet, nil)

	isOwner, err := tweetService.IsTweetOwner(context.Background(), tweetID, userID)
	assert.NoError(t, err)
	assert.True(t, isOwner)
	mockStorage.AssertExpectations(t)
}

func TestTweetService_IncrementTweetViews(t *testing.T) {
	tweetService, mockStorage, _ := setupTweetService()
	userID := "user1"
	tweetID := "1"

	mockStorage.On("Tweet").Return(mockStorage)
	mockStorage.On("IncrementTweetViews", mock.Anything, userID, tweetID).Return(nil)

	err := tweetService.IncrementTweetViews(context.Background(), userID, tweetID)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}
