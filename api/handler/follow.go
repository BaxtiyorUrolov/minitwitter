package handler

import (
	"context"
	"net/http"
	"twitter/api/models"

	"github.com/gin-gonic/gin"
)

// FollowUser godoc
// @Router       /api/v1/user/{id}/follow [POST]
// @Summary      Follow a user
// @Description  Follow a user by the logged-in user
// @Tags         follow
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID to follow"
// @Success      200  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Failure      403  {object}  models.Response
func (h *Handler) FollowUser(c *gin.Context) {
	followingID := c.Param("id")
	followerID := c.GetString("user_id")

	follow := models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	err := h.services.Follow().FollowUser(context.Background(), follow)
	if err != nil {
		handleResponse(c, h.log, "Error while following user", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "User followed successfully", http.StatusOK, nil)
}

// UnfollowUser godoc
// @Router       /api/v1/user/{id}/unfollow [DELETE]
// @Summary      Unfollow a user
// @Description  Unfollow a user by the logged-in user
// @Tags         follow
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID to unfollow"
// @Success      200  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Failure      403  {object}  models.Response
func (h *Handler) UnfollowUser(c *gin.Context) {
	followingID := c.Param("id")
	followerID := c.GetString("user_id")

	err := h.services.Follow().UnfollowUser(context.Background(), followerID, followingID)
	if err != nil {
		handleResponse(c, h.log, "Error while unfollowing user", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "User unfollowed successfully", http.StatusOK, nil)
}

// GetFollowers godoc
// @Router       /api/v1/user/{id}/followers [GET]
// @Summary      Get followers of a user
// @Description  Get the followers of a specific user
// @Tags         follow
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) GetFollowers(c *gin.Context) {
	userID := c.Param("id")

	followers, err := h.services.Follow().GetFollowers(context.Background(), userID)
	if err != nil {
		handleResponse(c, h.log, "Error while getting followers", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Success", http.StatusOK, followers)
}

// GetFollowings godoc
// @Router       /api/v1/user/{id}/followings [GET]
// @Summary      Get followings of a user
// @Description  Get the followings of a specific user
// @Tags         follow
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) GetFollowings(c *gin.Context) {
	userID := c.Param("id")

	followings, err := h.services.Follow().GetFollowings(context.Background(), userID)
	if err != nil {
		handleResponse(c, h.log, "Error while getting followings", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Success", http.StatusOK, followings)
}
