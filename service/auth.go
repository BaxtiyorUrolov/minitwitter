package service

import (
	"context"
	"twitter/api/models"
	"twitter/pkg/jwt"
	"twitter/pkg/logger"
	"twitter/pkg/security"
	"twitter/storage"
)

type authService struct {
	storage storage.IStorage
	log     logger.Logger
}

func NewAuthService(storage storage.IStorage, log logger.Logger) authService {
	return authService{
		storage: storage,
		log:     log,
	}
}

func (a authService) Login(ctx context.Context, loginRequest models.LoginRequest) (models.LoginResponse, error) {

	user, err := a.storage.User().GetUserCredentials(ctx, loginRequest.Login)
	if err != nil {
		return models.LoginResponse{}, err
	}

	if err = security.CompareHashAndPassword(user.Password, loginRequest.Password); err != nil {
		return models.LoginResponse{}, err
	}

	m := make(map[string]interface{})
	m["user_id"] = user.ID

	accessToken, refreshToken, err := jwt.GenerateJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for customer login", logger.Error(err))
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

	return models.LoginResponse{}, nil
}
