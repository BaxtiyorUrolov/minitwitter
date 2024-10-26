package service

import (
	"context"
	"encoding/json"
	"fmt"
	"twitter/api/models"
	"twitter/pkg/kafka"
	"twitter/pkg/logger"
	"twitter/storage"
)

type likeService struct {
	storage  storage.IStorage
	log      logger.Logger
	producer *kafka.Producer
}

func NewLikeService(storage storage.IStorage, log logger.Logger, producer *kafka.Producer) likeService {
	return likeService{
		storage:  storage,
		log:      log,
		producer: producer,
	}
}

type LikeNotificationEvent struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (l likeService) LikeTweet(ctx context.Context, like models.Like) error {
	l.log.Info("Liking tweet service layer", logger.Any("like", like))

	tweet, err := l.storage.Tweet().GetTweet(ctx, like.TweetID) // Tweet ID bilan tweetni oling
	if err != nil {
		l.log.Error("Error while getting tweet in like service", logger.Error(err))
		return err
	}

	user, err := l.storage.User().GetByID(ctx, models.PrimaryKey{
		ID: tweet.UserID,
	})
	if err != nil {
		l.log.Error("Error while getting tweet owner", logger.Error(err))
		return err
	}

	// Tweetni yoqtirish
	err = l.storage.Like().LikeTweet(ctx, like)
	if err != nil {
		l.log.Error("Error while liking tweet", logger.Error(err))
		return err
	}

	// Kafka xabarini yuborish
	notification := LikeNotificationEvent{
		Email:   user.Email,
		Message: fmt.Sprintf("User %s liked your tweet!", like.UserID),
	}

	notificationData, _ := json.Marshal(notification)
	if err := l.producer.SendMessage(string(notificationData)); err != nil {
		l.log.Error("Error while sending message to Kafka", logger.Error(err))
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
