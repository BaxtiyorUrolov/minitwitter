package handler

import (
	"context"
	"net/http"
	"twitter/api/models"

	"github.com/gin-gonic/gin"
)

// LikeTweet godoc
// @Router       /api/v1/tweet/{id}/like [POST]
// @Summary      Like a tweet
// @Description  Like a tweet by the logged-in user
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Tweet ID"
// @Success      200  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Failure      403  {object}  models.Response
func (h *Handler) LikeTweet(c *gin.Context) {
	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	like := models.Like{
		UserID:  userID,
		TweetID: tweetID,
	}

	err := h.services.Like().LikeTweet(context.Background(), like)
	if err != nil {
		handleResponse(c, h.log, "Error while liking tweet", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Tweet liked successfully", http.StatusOK, nil)
}

// UnlikeTweet godoc
// @Router       /api/v1/tweet/{id}/unlike [DELETE]
// @Summary      Unlike a tweet
// @Description  Unlike a tweet by the logged-in user
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Tweet ID"
// @Success      200  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Failure      403  {object}  models.Response
func (h *Handler) UnlikeTweet(c *gin.Context) {
	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.services.Like().UnlikeTweet(context.Background(), userID, tweetID)
	if err != nil {
		handleResponse(c, h.log, "Error while unliking tweet", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Tweet unliked successfully", http.StatusOK, nil)
}

// GetLikeCount godoc
// @Router       /api/v1/tweet/{id}/like-count [GET]
// @Summary      Get like count for a tweet
// @Description  Get the number of likes for a specific tweet
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tweet ID"
// @Success      200  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) GetLikeCount(c *gin.Context) {
	tweetID := c.Param("id")

	count, err := h.services.Like().GetLikeCount(context.Background(), tweetID)
	if err != nil {
		handleResponse(c, h.log, "Error while getting like count", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Success", http.StatusOK, gin.H{"like_count": count})
}
