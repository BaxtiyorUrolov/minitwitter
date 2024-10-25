package postgres

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"twitter/api/models"
)

func TestTweetRepo_CreateTweet(t *testing.T) {

	repo, ctx := setup(t)

	// Foydalanuvchini yaratish
	userRepo := repo.User() // Agar `UserRepo` mavjud bo'lsa
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	// Tweet yaratish
	tweet := models.CreateTweet{
		UserID:  userID, // To'g'ri user ID ishlatiladi
		Content: "This is a test tweet",
		Media:   "test_media_url",
	}

	_, err = repo.Tweet().CreateTweet(ctx, tweet)
	assert.NoError(t, err, "CreateTweet should not return an error")
}

func TestTweetRepo_GetTweet(t *testing.T) {
	repo, ctx := setup(t)

	userRepo := repo.User() // Agar `UserRepo` mavjud bo'lsa
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	// Create a tweet
	tweet := models.CreateTweet{
		UserID:  userID,
		Content: "This is a test tweet",
		Media:   "test_media_url",
	}
	id, err := repo.Tweet().CreateTweet(ctx, tweet)
	assert.NoError(t, err)

	// Get the tweet by ID
	fetchedTweet, err := repo.Tweet().GetTweet(ctx, id) // Tweet ID ni oling
	assert.NoError(t, err, "GetTweet should not return an error")
	assert.Equal(t, tweet.Content, fetchedTweet.Content, "Content should match")
	assert.Equal(t, tweet.Media, fetchedTweet.Media, "Media should match")
}

func TestTweetRepo_UpdateTweet(t *testing.T) {
	repo, ctx := setup(t)

	userRepo := repo.User() // Agar `UserRepo` mavjud bo'lsa
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	// Create a tweet
	tweet := models.CreateTweet{
		UserID:  userID,
		Content: "This is a test tweet",
		Media:   "test_media_url",
	}
	id, err := repo.Tweet().CreateTweet(ctx, tweet)
	assert.NoError(t, err)

	// Update the tweet
	update := models.UpdateTweet{
		ID:      id,
		Content: "Updated content",
		Media:   "updated_media_url",
	}
	err = repo.Tweet().UpdateTweet(ctx, update)
	assert.NoError(t, err, "UpdateTweet should not return an error")

	// Verify the update
	updatedTweet, err := repo.Tweet().GetTweet(ctx, id)
	assert.NoError(t, err, "GetTweet should not return an error after update")
	assert.Equal(t, update.Content, updatedTweet.Content, "Content should match the updated content")
	assert.Equal(t, update.Media, updatedTweet.Media, "Media should match the updated media")
}

func TestTweetRepo_DeleteTweet(t *testing.T) {
	repo, ctx := setup(t)

	userRepo := repo.User() // Agar `UserRepo` mavjud bo'lsa
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	// Create a tweet
	tweet := models.CreateTweet{
		UserID:  userID,
		Content: "This is a test tweet",
		Media:   "test_media_url",
	}
	id, err := repo.Tweet().CreateTweet(ctx, tweet)
	assert.NoError(t, err)

	// Delete the tweet
	err = repo.Tweet().DeleteTweet(ctx, id)
	assert.NoError(t, err, "DeleteTweet should not return an error")

	// Verify that the tweet is deleted
	_, err = repo.Tweet().GetTweet(ctx, id)
	assert.Error(t, err, "GetTweet for deleted tweet should return an error")
}

func TestTweetRepo_GetTweetList(t *testing.T) {
	repo, ctx := setup(t)

	// Foydalanuvchini yaratish
	userRepo := repo.User()
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	_, err = repo.Tweet().CreateTweet(ctx, models.CreateTweet{
		UserID:  userID,
		Content: "BMC",
		Media:   "link",
	})

	_, err = repo.Tweet().CreateTweet(ctx, models.CreateTweet{
		UserID:  userID,
		Content: "GrScan",
		Media:   "link",
	})

	assert.NoError(t, err)

	// Fetch the tweet list
	req := models.GetListRequest{
		Page:   1,
		Limit:  3,
		Search: "Test",
	}
	resp, err := repo.Tweet().GetTweetList(ctx, req)
	assert.NoError(t, err, "GetTweetList should not return an error")
	assert.Len(t, resp.Tweets, 3, "There should be 3 tweets returned")
}

func TestTweetRepo_IncrementTweetViews(t *testing.T) {
	repo, ctx := setup(t)

	// Foydalanuvchi yaratish
	userRepo := repo.User()
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := userRepo.Create(ctx, createUser)
	assert.NoError(t, err)

	// Yangi foydalanuvchini yaratish, bu foydalanuvchi tvitni ko'radi
	viewingUser := models.CreateUser{
		Name:     "Viewing User",
		UserName: "viewinguser",
		Email:    "viewinguser@example.com",
		Password: "password123",
	}
	viewingUserID, err := userRepo.Create(ctx, viewingUser)
	assert.NoError(t, err)

	// Tvit yaratish
	tweet := models.CreateTweet{
		UserID:  userID,
		Content: "This is a test tweet",
		Media:   "test_media_url",
	}
	id, err := repo.Tweet().CreateTweet(ctx, tweet)
	assert.NoError(t, err)

	// Tvit uchun ko'rishlar sonini oshirish
	err = repo.Tweet().IncrementTweetViews(ctx, viewingUserID, id) // To'g'ri UUID ishlatilmoqda
	assert.NoError(t, err, "IncrementTweetViews should not return an error")

	// Ko'rishlar soni oshganini tekshiring
	fetchedTweet, err := repo.Tweet().GetTweet(ctx, id)
	assert.NoError(t, err, "GetTweet should not return an error after views increment")
	assert.Equal(t, 1, fetchedTweet.ViewsCount, "Views count should be incremented")
}
