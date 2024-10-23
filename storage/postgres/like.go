package postgres

import (
	"context"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type likeRepo struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewLikeRepo(db *pgxpool.Pool, log logger.Logger) storage.ILikeStorage {
	return &likeRepo{
		db:  db,
		log: log,
	}
}

func (l *likeRepo) LikeTweet(ctx context.Context, like models.Like) error {
	_, err := l.db.Exec(ctx, `
		INSERT INTO likes (user_id, tweet_id, created_at) 
		VALUES ($1, $2, NOW())
		`,
		like.UserID,
		like.TweetID,
	)
	if err != nil {
		l.log.Error("error while liking tweet", logger.Error(err))
		return err
	}
	return nil
}

func (l *likeRepo) UnlikeTweet(ctx context.Context, userID, tweetID string) error {
	_, err := l.db.Exec(ctx, `DELETE FROM likes WHERE user_id = $1 AND tweet_id = $2`, userID, tweetID)
	if err != nil {
		l.log.Error("error while unliking tweet", logger.Error(err))
		return err
	}
	return nil
}

func (l *likeRepo) GetLikeCount(ctx context.Context, tweetID string) (int, error) {
	var likeCount int
	err := l.db.QueryRow(ctx, `SELECT COUNT(*) FROM likes WHERE tweet_id = $1`, tweetID).Scan(&likeCount)
	if err != nil {
		l.log.Error("error while fetching like count", logger.Error(err))
		return 0, err
	}
	return likeCount, nil
}
