package service

import (
	"context"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/storage"
)

type likeService struct {
	storage storage.IStorage
	log     logger.Logger
}

func NewLikeService(storage storage.IStorage, log logger.Logger) likeService {
	return likeService{
		storage: storage,
		log:     log,
	}
}

func (l likeService) LikeTweet(ctx context.Context, like models.Like) error {
	l.log.Info("Liking tweet service layer", logger.Any("like", like))

	err := l.storage.Like().LikeTweet(ctx, like)
	if err != nil {
		l.log.Error("Error while liking tweet", logger.Error(err))
		return err
	}

	return nil
}

func (l likeService) UnlikeTweet(ctx context.Context, userID, tweetID string) error {
	l.log.Info("Unliking tweet service layer", logger.Any("user_id", userID), logger.Any("tweet_id", tweetID))

	err := l.storage.Like().UnlikeTweet(ctx, userID, tweetID)
	if err != nil {
		l.log.Error("Error while unliking tweet", logger.Error(err))
		return err
	}

	return nil
}

func (l likeService) GetLikeCount(ctx context.Context, tweetID string) (int, error) {
	l.log.Info("Getting like count for tweet in service layer", logger.Any("tweet_id", tweetID))

	count, err := l.storage.Like().GetLikeCount(ctx, tweetID)
	if err != nil {
		l.log.Error("Error in service layer when getting like count", logger.Error(err))
		return 0, err
	}

	return count, nil
}
