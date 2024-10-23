package service

import (
	"context"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/pkg/security"
	"twitter/storage"
)

type userService struct {
	storage storage.IStorage
	log     logger.Logger
}

func NewUserService(storage storage.IStorage, log logger.Logger) userService {
	return userService{
		storage: storage,
		log:     log,
	}
}

func (u userService) Create(ctx context.Context, createUser models.CreateUser) (models.User, error) {
	u.log.Info("user create service layer", logger.Any("user", createUser))

	// Parolni hash qilish
	password, err := security.HashPassword(createUser.Password)
	if err != nil {
		u.log.Error("Error while hashing password", logger.Error(err))
		return models.User{}, err
	}
	createUser.Password = password

	pKey, err := u.storage.User().Create(ctx, createUser)
	if err != nil {
		u.log.Error("Error while creating user", logger.Error(err))
		return models.User{}, err
	}

	user, err := u.storage.User().GetByID(ctx, models.PrimaryKey{ID: pKey})
	if err != nil {
		u.log.Error("Error in service layer when getting user by id", logger.Error(err))
		return user, err
	}

	return user, nil
}

func (u userService) IsLoginExist(ctx context.Context, login string) (bool, error) {
	exists, err := u.storage.User().IsUserNameExist(ctx, login)
	if err != nil {
		u.log.Error("Error in service layer when checking user", logger.Error(err))
		return false, err
	}

	return exists, nil
}

func (u userService) GetByID(ctx context.Context, id models.PrimaryKey) (models.User, error) {
	user, err := u.storage.User().GetByID(ctx, id)
	if err != nil {
		u.log.Error("Error in service layer when getting user by id", logger.Error(err))
		return models.User{}, err
	}

	return user, nil
}

func (u userService) Update(ctx context.Context, updateUser models.UpdateUser) error {
	u.log.Info("user update service layer", logger.Any("user", updateUser))

	err := u.storage.User().Update(ctx, updateUser)
	if err != nil {
		u.log.Error("Error while updating user", logger.Error(err))
		return err
	}

	return nil
}

func (u userService) Delete(ctx context.Context, id models.PrimaryKey) error {

	err := u.storage.User().Delete(ctx, id)
	if err != nil {
		u.log.Error("Error while deleting user", logger.Error(err))
		return err
	}

	return nil
}

func (u userService) GetAllUsers(ctx context.Context, request models.GetListRequest) (models.UsersResponse, error) {
	users, err := u.storage.User().GetList(ctx, request)
	if err != nil {
		u.log.Error("Error in service layer when getting all users", logger.Error(err))
		return models.UsersResponse{}, err
	}

	return users, nil
}
