package service

import (
	"context"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/storage"
)

type followService struct {
	storage storage.IStorage
	log     logger.Logger
}

func NewFollowService(storage storage.IStorage, log logger.Logger) followService {
	return followService{
		storage: storage,
		log:     log,
	}
}

func (f followService) FollowUser(ctx context.Context, follow models.Follow) error {
	f.log.Info("Following user service layer", logger.Any("follow", follow))

	err := f.storage.Follow().FollowUser(ctx, follow)
	if err != nil {
		f.log.Error("Error while following user", logger.Error(err))
		return err
	}

	return nil
}

func (f followService) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	f.log.Info("Unfollowing user service layer", logger.Any("follower_id", followerID), logger.Any("following_id", followingID))

	err := f.storage.Follow().UnfollowUser(ctx, followerID, followingID)
	if err != nil {
		f.log.Error("Error while unfollowing user", logger.Error(err))
		return err
	}

	return nil
}

func (f followService) GetFollowers(ctx context.Context, userID string) ([]models.Follow, error) {
	f.log.Info("Getting followers in service layer", logger.Any("user_id", userID))

	followers, err := f.storage.Follow().GetFollowers(ctx, userID)
	if err != nil {
		f.log.Error("Error in service layer when getting followers", logger.Error(err))
		return nil, err
	}

	return followers, nil
}

func (f followService) GetFollowings(ctx context.Context, userID string) ([]models.Follow, error) {
	f.log.Info("Getting followings in service layer", logger.Any("user_id", userID))

	followings, err := f.storage.Follow().GetFollowing(ctx, userID)
	if err != nil {
		f.log.Error("Error in service layer when getting followings", logger.Error(err))
		return nil, err
	}

	return followings, nil
}
