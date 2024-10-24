package service

import (
	"context"
	"errors"
	"fmt"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/storage"
)

type tweetService struct {
	storage storage.IStorage
	log     logger.Logger
}

func NewTweetService(storage storage.IStorage, log logger.Logger) tweetService {
	return tweetService{
		storage: storage,
		log:     log,
	}
}

func (t tweetService) Create(ctx context.Context, createTweet models.CreateTweet) error {

	err := t.storage.Tweet().CreateTweet(ctx, createTweet)
	if err != nil {
		t.log.Error("Error while creating tweet", logger.Error(err))
		return err
	}

	return nil
}

func (t tweetService) Get(ctx context.Context, tweetID string) (models.Tweet, error) {
	tweet, err := t.storage.Tweet().GetTweet(ctx, tweetID)
	if err != nil {
		return tweet, err
	}

	return tweet, nil
}

func (t tweetService) Update(ctx context.Context, tweet models.UpdateTweet) error {

	err := t.storage.Tweet().UpdateTweet(ctx, tweet)
	if err != nil {
		t.log.Error("Error while updating tweet", logger.Error(err))
		return err
	}

	return nil
}

func (t tweetService) Delete(ctx context.Context, tweetID string) error {

	err := t.storage.Tweet().DeleteTweet(ctx, tweetID)
	if err != nil {
		t.log.Error("Error while deleting tweet", logger.Error(err))
		return err
	}

	return nil
}

func (t tweetService) GetTweetsByUser(ctx context.Context, userID string) (models.TweetsResponse, error) {

	tweets, err := t.storage.Tweet().ListTweetsByUser(ctx, userID)
	if err != nil {
		t.log.Error("Error in service layer when listing tweets by user", logger.Error(err))
		return models.TweetsResponse{}, err
	}

	fmt.Println(tweets.Count)

	return tweets, nil
}

func (t tweetService) GetList(ctx context.Context, request models.GetListRequest) (models.TweetsResponse, error) {

	tweets, err := t.storage.Tweet().GetTweetList(ctx, request)
	if err != nil {
		t.log.Error("Error in service layer when getting tweets", logger.Error(err))
		return models.TweetsResponse{}, err
	}

	return tweets, nil
}

func (t tweetService) IsTweetOwner(ctx context.Context, tweetID string, userID string) (bool, error) {
	fmt.Println("id: ", tweetID)
	tweet, err := t.storage.Tweet().GetTweet(ctx, tweetID)
	if err != nil {
		return false, err
	}

	if tweet.UserID != userID {
		return false, errors.New("user is not the owner of this tweet")
	}

	return true, nil
}

func (t tweetService) IncrementTweetViews(ctx context.Context, userID string, tweetID string) error {

	err := t.storage.Tweet().IncrementTweetViews(ctx, userID, tweetID)
	if err != nil {
		t.log.Error("Error in service layer when incrementing tweet views", logger.Error(err))
		return err
	}
	return nil
}
