package postgres

import (
	"context"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type followRepo struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewFollowRepo(db *pgxpool.Pool, log logger.Logger) storage.IFollowStorage {
	return &followRepo{
		db:  db,
		log: log,
	}
}

func (f *followRepo) FollowUser(ctx context.Context, follow models.Follow) error {
	_, err := f.db.Exec(ctx, `
		INSERT INTO follows (follower_id, following_id, created_at) 
		VALUES ($1, $2, NOW())
		`,
		follow.FollowerID,
		follow.FollowingID,
	)
	if err != nil {
		f.log.Error("error while following user", logger.Error(err))
		return err
	}
	return nil
}

func (f *followRepo) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	_, err := f.db.Exec(ctx, `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`, followerID, followingID)
	if err != nil {
		f.log.Error("error while unfollowing user", logger.Error(err))
		return err
	}
	return nil
}

func (f *followRepo) GetFollowers(ctx context.Context, userID string) ([]models.Follow, error) {
	rows, err := f.db.Query(ctx, `SELECT follower_id, following_id, created_at FROM follows WHERE following_id = $1`, userID)
	if err != nil {
		f.log.Error("error while fetching followers", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	var followers []models.Follow
	for rows.Next() {
		follow := models.Follow{}
		err = rows.Scan(&follow.FollowerID, &follow.FollowingID, &follow.CreatedAt)
		if err != nil {
			f.log.Error("error while scanning followers", logger.Error(err))
			continue
		}
		followers = append(followers, follow)
	}
	return followers, nil
}

func (f *followRepo) GetFollowing(ctx context.Context, userID string) ([]models.Follow, error) {
	rows, err := f.db.Query(ctx, `SELECT follower_id, following_id, created_at FROM follows WHERE follower_id = $1`, userID)
	if err != nil {
		f.log.Error("error while fetching following", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	var followings []models.Follow
	for rows.Next() {
		follow := models.Follow{}
		err = rows.Scan(&follow.FollowerID, &follow.FollowingID, &follow.CreatedAt)
		if err != nil {
			f.log.Error("error while scanning following", logger.Error(err))
			continue
		}
		followings = append(followings, follow)
	}
	return followings, nil
}
