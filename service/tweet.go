package service

import (
	"context"
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

func (t tweetService) Create(ctx context.Context, createTweet models.Tweet) (models.Tweet, error) {

	tweetID, err := t.storage.Tweet().CreateTweet(ctx, createTweet)
	if err != nil {
		t.log.Error("Error while creating tweet", logger.Error(err))
		return models.Tweet{}, err
	}

	tweet, err := t.storage.Tweet().GetTweet(ctx, tweetID)
	if err != nil {
		t.log.Error("Error in service layer when getting tweet by id", logger.Error(err))
		return tweet, err
	}

	return tweet, nil
}

func (t tweetService) Get(ctx context.Context, tweetID string) (models.Tweet, error) {
	tweet, err := t.storage.Tweet().GetTweet(ctx, tweetID)
	if err != nil {
		return tweet, err
	}

	return tweet, nil
}

func (t tweetService) Update(ctx context.Context, tweet models.Tweet) error {

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

func (t tweetService) ListTweetsByUser(ctx context.Context, userID string) ([]models.Tweet, error) {

	tweets, err := t.storage.Tweet().ListTweetsByUser(ctx, userID)
	if err != nil {
		t.log.Error("Error in service layer when listing tweets by user", logger.Error(err))
		return nil, err
	}

	return tweets, nil
}
