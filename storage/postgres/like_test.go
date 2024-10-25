package postgres

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"twitter/api/models"
)

func TestLikeRepo_LikeTweet(t *testing.T) {

	repo, ctx := setup(t)

	// Foydalanuvchi yaratish
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := repo.User().Create(ctx, createUser)
	assert.NoError(t, err)

	// Tvit yaratish
	createTweet := models.CreateTweet{
		UserID:  userID,
		Content: "This is a test tweet",
		Media:   "test_media_url",
	}
	tweetID, err := repo.Tweet().CreateTweet(ctx, createTweet)
	assert.NoError(t, err)

	// Tvitni like qilish
	like := models.Like{
		UserID:  userID,
		TweetID: tweetID,
	}
	err = repo.Like().LikeTweet(ctx, like)
	assert.NoError(t, err, "LikeTweet should not return an error")
}

func TestLikeRepo_UnlikeTweet(t *testing.T) {
	repo, ctx := setup(t)

	// Foydalanuvchi yaratish
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := repo.User().Create(ctx, createUser)
	assert.NoError(t, err)

	// Tvit yaratish
	createTweet := models.CreateTweet{
		UserID:  userID,
		Content: "This is a test tweet",
		Media:   "test_media_url",
	}
	tweetID, err := repo.Tweet().CreateTweet(ctx, createTweet)
	assert.NoError(t, err)

	// Tvitni like qilish va keyin unlike qilish
	like := models.Like{
		UserID:  userID,
		TweetID: tweetID,
	}
	err = repo.Like().LikeTweet(ctx, like)
	assert.NoError(t, err, "LikeTweet should not return an error")

	err = repo.Like().UnlikeTweet(ctx, userID, tweetID)
	assert.NoError(t, err, "UnlikeTweet should not return an error")

	// Tvitning like holatini tekshirish
	likeCount, err := repo.Like().GetLikeCount(ctx, tweetID)
	assert.NoError(t, err)
	assert.Equal(t, 0, likeCount, "Like count should be 0 after unliking")
}

func TestLikeRepo_GetLikeCount(t *testing.T) {
	repo, ctx := setup(t)

	// Foydalanuvchi yaratish
	createUser := models.CreateUser{
		Name:     "Test User",
		UserName: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userID, err := repo.User().Create(ctx, createUser)
	assert.NoError(t, err)

	// Tvit yaratish
	createTweet := models.CreateTweet{
		UserID:  userID,
		Content: "This is a test tweet",
		Media:   "test_media_url",
	}
	tweetID, err := repo.Tweet().CreateTweet(ctx, createTweet)
	assert.NoError(t, err)

	// Tvitni like qilish
	like := models.Like{
		UserID:  userID,
		TweetID: tweetID,
	}
	err = repo.Like().LikeTweet(ctx, like)
	assert.NoError(t, err, "LikeTweet should not return an error")

	// Like count-ni tekshirish
	likeCount, err := repo.Like().GetLikeCount(ctx, tweetID)
	assert.NoError(t, err)
	assert.Equal(t, 1, likeCount, "Like count should be 1 after liking")
}
