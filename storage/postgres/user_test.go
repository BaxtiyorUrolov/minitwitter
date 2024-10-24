package postgres

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"twitter/api/models"
	"twitter/config"
	"twitter/pkg/logger"
	"twitter/storage"
)

func setup(t *testing.T) (storage.IStorage, logger.Logger, context.Context) {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	// Run migrations
	runMigrations(t, cfg, log)

	// Initialize store
	store, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Fatalf("Error while connecting to DB: %v", err)
	}

	return store, log, context.Background()
}

// Helper function to run migrations
func runMigrations(t *testing.T, cfg config.Config, log logger.Logger) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB,
	)

	// Migrate up
	m, err := migrate.New("file://migrations/postgres", url)
	if err != nil {
		t.Fatalf("error while setting up migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatalf("error while running migration: %v", err)
	}
}

func TestUserRepo_Create(t *testing.T) {
	pgStore, _, ctx := setup(t)

	// UserRepo dan foydalanamiz
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

	// Create a user for testing
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	// Test ma'lumotlar bazasida mavjud bo'lgan foydalanuvchini oling
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

	// Create a user for testing
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

	// Verify that the user is deleted
	_, err = userRepo.GetByID(ctx, pKey)
	assert.Error(t, err, "GetByID for deleted user should return an error")
}

func TestUserRepo_IsUserNameExist(t *testing.T) {
	pgStore, _, ctx := setup(t)

	userRepo := pgStore.User()

	// Create a user for testing
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	_, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	// Check if the username exists
	exists, err := userRepo.IsUserNameExist(ctx, "testuser")
	assert.NoError(t, err, "IsUserNameExist should not return an error")
	assert.True(t, exists, "Username should exist")
}
