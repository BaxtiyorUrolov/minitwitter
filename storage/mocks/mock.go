package mocks

import (
	"context"
	"twitter/api/models"
	"twitter/storage"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Close() {
	m.Called()
}

func (m *MockStorage) User() storage.IUserStorage {
	args := m.Called()
	return args.Get(0).(storage.IUserStorage)
}

func (m *MockStorage) Tweet() storage.ITweetStorage {
	args := m.Called()
	return args.Get(0).(storage.ITweetStorage)
}

func (m *MockStorage) Like() storage.ILikeStorage {
	args := m.Called()
	return args.Get(0).(storage.ILikeStorage)
}

func (m *MockStorage) Follow() storage.IFollowStorage {
	args := m.Called()
	return args.Get(0).(storage.IFollowStorage)
}

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) Create(ctx context.Context, createUser models.CreateUser) (string, error) {
	args := m.Called(ctx, createUser)
	return args.String(0), args.Error(1)
}

func (m *MockUserStorage) GetByID(ctx context.Context, id models.PrimaryKey) (models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserStorage) IsUserNameExist(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserStorage) Update(ctx context.Context, updateUser models.UpdateUser) error {
	args := m.Called(ctx, updateUser)
	return args.Error(0)
}

func (m *MockUserStorage) Delete(ctx context.Context, id models.PrimaryKey) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserStorage) GetList(ctx context.Context, request models.GetListRequest) (models.UsersResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(models.UsersResponse), args.Error(1)
}

func (m *MockUserStorage) GetUserCredentials(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

// Tweet

func (m *MockStorage) CreateTweet(ctx context.Context, tweet models.CreateTweet) (string, error) {
	args := m.Called(ctx, tweet)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GetTweet(ctx context.Context, tweetID string) (models.Tweet, error) {
	args := m.Called(ctx, tweetID)
	return args.Get(0).(models.Tweet), args.Error(1)
}

func (m *MockStorage) UpdateTweet(ctx context.Context, tweet models.UpdateTweet) error {
	args := m.Called(ctx, tweet)
	return args.Error(0)
}

func (m *MockStorage) DeleteTweet(ctx context.Context, tweetID string) error {
	args := m.Called(ctx, tweetID)
	return args.Error(0)
}

func (m *MockStorage) ListTweetsByUser(ctx context.Context, userID string) (models.TweetsResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.TweetsResponse), args.Error(1)
}

func (m *MockStorage) GetTweetList(ctx context.Context, request models.GetListRequest) (models.TweetsResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(models.TweetsResponse), args.Error(1)
}

func (m *MockStorage) IncrementTweetViews(ctx context.Context, tweetID string, userID string) error {
	args := m.Called(ctx, tweetID, userID)
	return args.Error(0)
}

//Like storage

func (m *MockStorage) LikeTweet(ctx context.Context, like models.Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *MockStorage) UnlikeTweet(ctx context.Context, userID, tweetID string) error {
	args := m.Called(ctx, userID, tweetID)
	return args.Error(0)
}

func (m *MockStorage) GetLikeCount(ctx context.Context, tweetID string) (int, error) {
	args := m.Called(ctx, tweetID)
	return args.Int(0), args.Error(1)
}

// Follow storage

func (m *MockStorage) FollowUser(ctx context.Context, follow models.Follow) error {
	args := m.Called(ctx, follow)
	return args.Error(0)
}

func (m *MockStorage) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	args := m.Called(ctx, followerID, followingID)
	return args.Error(0)
}

func (m *MockStorage) GetFollowers(ctx context.Context, userID string) ([]models.Follow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Follow), args.Error(1)
}

func (m *MockStorage) GetFollowing(ctx context.Context, userID string) ([]models.Follow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Follow), args.Error(1)
}
