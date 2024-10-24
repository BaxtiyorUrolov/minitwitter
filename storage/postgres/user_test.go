package postgres

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"twitter/api/models"
	"twitter/config"
	"twitter/pkg/logger"
	"twitter/storage"
)

func setup(t *testing.T) (storage.IStorage, logger.Logger, context.Context) {

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
	}
	fmt.Println("Current directory:", dir)

	err = godotenv.Load(".env")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	// Load config and logger
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	// Run migrations
	store, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Fatalf("Error while connecting to DB: %v", err)
	}

	return store, log, context.Background()
}

func TestUserRepo_Create(t *testing.T) {

	err := os.Chdir("../../")
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}

	pgStore, _, ctx := setup(t)

	userRepo := pgStore.User()

	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err, "Create should not return an error")
	assert.NotEmpty(t, userID, "User ID should not be empty")
}

func TestUserRepo_GetByID(t *testing.T) {

	pgStore, _, ctx := setup(t)

	userRepo := pgStore.User()

	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	pKey := models.PrimaryKey{
		ID: userID,
	}

	user, err := userRepo.GetByID(ctx, pKey)
	assert.NoError(t, err, "GetByID should not return an error")
	assert.Equal(t, pKey.ID, user.ID, "User ID should match")
}

func TestUserRepo_Update(t *testing.T) {

	pgStore, _, ctx := setup(t)

	userRepo := pgStore.User()

	// Create a user for testing
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	updateUser := models.UpdateUser{
		ID:             userID,
		Name:           "Updated Name",
		UserName:       "updatedusername",
		Email:          "updatedemail@example.com",
		Bio:            "Updated bio",
		ProfilePicture: "updated-profile-pic-url",
	}

	err = userRepo.Update(ctx, updateUser)
	assert.NoError(t, err, "Update should not return an error")

	// Verify the update
	pKey := models.PrimaryKey{
		ID: userID,
	}

	updatedUser, err := userRepo.GetByID(ctx, pKey)
	assert.NoError(t, err, "GetByID after update should not return an error")
	assert.Equal(t, updateUser.Name, updatedUser.Name, "User name should match the updated name")
	assert.Equal(t, updateUser.Email, updatedUser.Email, "User email should match the updated email")
}

func TestUserRepo_Delete(t *testing.T) {

	pgStore, _, ctx := setup(t)

	userRepo := pgStore.User()

	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	pKey := models.PrimaryKey{
		ID: userID,
	}

	err = userRepo.Delete(ctx, pKey)
	assert.NoError(t, err, "Delete should not return an error")

}

func TestUserRepo_IsUserNameExist(t *testing.T) {

	pgStore, _, ctx := setup(t)

	userRepo := pgStore.User()

	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	_, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	exists, err := userRepo.IsUserNameExist(ctx, "testuser")
	assert.NoError(t, err, "IsUserNameExist should not return an error")
	assert.True(t, exists, "Username should exist")
}
