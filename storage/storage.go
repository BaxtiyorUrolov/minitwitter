package storage

import (
	"context"
	"twitter/api/models"
)

type IStorage interface {
	Close()
	User() IUserStorage
	Tweet() ITweetStorage
	Like() ILikeStorage
	Follow() IFollowStorage
}

type IUserStorage interface {
	Create(context.Context, models.CreateUser) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.User, error)
	IsUserNameExist(context.Context, string) (bool, error)
	Update(context.Context, models.UpdateUser) error
	GetList(context.Context, models.GetListRequest) (models.UsersResponse, error)
	Delete(context.Context, models.PrimaryKey) error
	GetUserCredentials(context.Context, string) (models.User, error)
}

type ITweetStorage interface {
	CreateTweet(ctx context.Context, tweet models.CreateTweet) error
	GetTweet(ctx context.Context, tweetID string) (models.Tweet, error)
	UpdateTweet(ctx context.Context, tweet models.UpdateTweet) error
	DeleteTweet(ctx context.Context, tweetID string) error
	ListTweetsByUser(ctx context.Context, userID string) (models.TweetsResponse, error)
	GetTweetList(ctx context.Context, request models.GetListRequest) (models.TweetsResponse, error)
	IncrementTweetViews(ctx context.Context, tweetID string, userID string) error
}

type ILikeStorage interface {
	LikeTweet(ctx context.Context, like models.Like) error
	UnlikeTweet(ctx context.Context, userID, tweetID string) error
	GetLikeCount(ctx context.Context, tweetID string) (int, error)
}

type IFollowStorage interface {
	FollowUser(ctx context.Context, follow models.Follow) error
	UnfollowUser(ctx context.Context, followerID, followingID string) error
	GetFollowers(ctx context.Context, userID string) ([]models.Follow, error)
	GetFollowing(ctx context.Context, userID string) ([]models.Follow, error)
}
