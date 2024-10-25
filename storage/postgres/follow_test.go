package postgres

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"twitter/api/models"
)

func TestFollowRepo_FollowUser(t *testing.T) {

	err := os.Chdir("../../")
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}

	repo, ctx := setup(t)

	user1 := models.CreateUser{
		Name:     "User 1",
		UserName: "user1",
		Email:    "user1@example.com",
		Password: "password123",
	}
	user1ID, err := repo.User().Create(ctx, user1)
	assert.NoError(t, err)

	user2 := models.CreateUser{
		Name:     "User 2",
		UserName: "user2",
		Email:    "user2@example.com",
		Password: "password123",
	}
	user2ID, err := repo.User().Create(ctx, user2)
	assert.NoError(t, err)

	follow := models.Follow{
		FollowerID:  user1ID,
		FollowingID: user2ID,
	}
	err = repo.Follow().FollowUser(ctx, follow)
	assert.NoError(t, err, "FollowUser should not return an error")
}

func TestFollowRepo_UnfollowUser(t *testing.T) {
	repo, ctx := setup(t)

	user1 := models.CreateUser{
		Name:     "User 1",
		UserName: "user1",
		Email:    "user1@example.com",
		Password: "password123",
	}
	user1ID, err := repo.User().Create(ctx, user1)
	assert.NoError(t, err)

	user2 := models.CreateUser{
		Name:     "User 2",
		UserName: "user2",
		Email:    "user2@example.com",
		Password: "password123",
	}
	user2ID, err := repo.User().Create(ctx, user2)
	assert.NoError(t, err)

	follow := models.Follow{
		FollowerID:  user1ID,
		FollowingID: user2ID,
	}
	err = repo.Follow().FollowUser(ctx, follow)
	assert.NoError(t, err)

	err = repo.Follow().UnfollowUser(ctx, user1ID, user2ID)
	assert.NoError(t, err, "UnfollowUser should not return an error")
}

func TestFollowRepo_GetFollowers(t *testing.T) {
	repo, ctx := setup(t)

	user1 := models.CreateUser{
		Name:     "User 1",
		UserName: "user1",
		Email:    "user1@example.com",
		Password: "password123",
	}
	user1ID, err := repo.User().Create(ctx, user1)
	assert.NoError(t, err)

	user2 := models.CreateUser{
		Name:     "User 2",
		UserName: "user2",
		Email:    "user2@example.com",
		Password: "password123",
	}
	user2ID, err := repo.User().Create(ctx, user2)
	assert.NoError(t, err)

	follow := models.Follow{
		FollowerID:  user1ID,
		FollowingID: user2ID,
	}
	err = repo.Follow().FollowUser(ctx, follow)
	assert.NoError(t, err)

	followers, err := repo.Follow().GetFollowers(ctx, user2ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(followers), "User 2 should have 1 follower")
	assert.Equal(t, user1ID, followers[0].FollowerID, "User 1 should be the follower")
}

func TestFollowRepo_GetFollowing(t *testing.T) {
	repo, ctx := setup(t)

	user1 := models.CreateUser{
		Name:     "User 1",
		UserName: "user1",
		Email:    "user1@example.com",
		Password: "password123",
	}
	user1ID, err := repo.User().Create(ctx, user1)
	assert.NoError(t, err)

	user2 := models.CreateUser{
		Name:     "User 2",
		UserName: "user2",
		Email:    "user2@example.com",
		Password: "password123",
	}
	user2ID, err := repo.User().Create(ctx, user2)
	assert.NoError(t, err)

	follow := models.Follow{
		FollowerID:  user1ID,
		FollowingID: user2ID,
	}
	err = repo.Follow().FollowUser(ctx, follow)
	assert.NoError(t, err)

	following, err := repo.Follow().GetFollowing(ctx, user1ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(following), "User 1 should be following 1 user")
	assert.Equal(t, user2ID, following[0].FollowingID, "User 2 should be the followed user")
}
