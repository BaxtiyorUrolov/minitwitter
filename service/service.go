package service

import (
	"twitter/pkg/logger"
	"twitter/storage"
)

type IServiceManager interface {
	User() userService
	Auth() authService
	Tweet() tweetService
	Like() likeService
	Follow() followService
}

type Service struct {
	userService   userService
	authService   authService
	tweetService  tweetService
	likeService   likeService
	followService followService
}

func New(storage storage.IStorage, log logger.Logger) Service {
	services := Service{}

	services.userService = NewUserService(storage, log)
	services.authService = NewAuthService(storage, log)
	services.tweetService = NewTweetService(storage, log)
	services.likeService = NewLikeService(storage, log)
	services.followService = NewFollowService(storage, log)

	return services
}

func (s Service) User() userService {
	return s.userService
}

func (s Service) Auth() authService {
	return s.authService
}

func (s Service) Tweet() tweetService { return s.tweetService }

func (s Service) Like() likeService { return s.likeService }

func (s Service) Follow() followService {
	return s.followService

}
