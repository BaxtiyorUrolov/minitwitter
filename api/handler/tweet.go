package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"twitter/api/models"
)

// CreateTweet godoc
// @Router       /api/v1/tweet [POST]
// @Summary      Creates a new tweet
// @Description  Creates a new tweet
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        tweet body models.CreateTweet false "tweet"
// @Success      201  {object}  models.Tweet
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CreateTweet(c *gin.Context) {
	tweet := models.CreateTweet{}
	if err := c.ShouldBindJSON(&tweet); err != nil {
		handleResponse(c, h.log, "Error while binding request body", http.StatusBadRequest, err)
		return
	}

	userID, _ := c.Get("user_id")
	tweet.UserID = userID.(string)

	err := h.services.Tweet().Create(context.Background(), tweet)
	if err != nil {
		handleResponse(c, h.log, "Error while creating tweet", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Tweet created successfully", http.StatusCreated, tweet)
}

// GetTweet godoc
// @Router       /api/v1/tweet/{id} [GET]
// @Summary      Get a tweet by ID
// @Description  Get a tweet by ID
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Param        id path string true "Tweet ID"
// @Success      200  {object}  models.Tweet
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) GetTweet(c *gin.Context) {
	tweetID := c.Param("id")

	tweet, err := h.services.Tweet().Get(context.Background(), tweetID)
	if err != nil {
		handleResponse(c, h.log, "Tweet not found", http.StatusNotFound, err)
		return
	}

	handleResponse(c, h.log, "Success", http.StatusOK, tweet)
}

// UpdateTweet godoc
// @Router       /api/v1/tweet/{id} [PUT]
// @Summary      Update a tweet
// @Description  Update a tweet
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        id path string true "Tweet ID"
// @Param        tweet body models.UpdateTweet false "tweet"
// @Success      200  {object}  models.Tweet
// @Failure      400  {object}  models.Response
// @Failure      403  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) UpdateTweet(c *gin.Context) {
	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	// Check if the tweet belongs to the authenticated user
	isOwner, err := h.services.Tweet().IsTweetOwner(context.Background(), tweetID, userID)
	if err != nil {
		handleResponse(c, h.log, "Tweet not found or user is not the owner", http.StatusForbidden, err)
		return
	}

	if !isOwner {
		handleResponse(c, h.log, "Unauthorized to update this tweet", http.StatusForbidden, nil)
		return
	}

	tweet := models.UpdateTweet{}
	if err := c.ShouldBindJSON(&tweet); err != nil {
		handleResponse(c, h.log, "Error while binding request body", http.StatusBadRequest, err)
		return
	}

	tweet.UserID = userID
	tweet.ID = tweetID

	err = h.services.Tweet().Update(context.Background(), tweet)
	if err != nil {
		handleResponse(c, h.log, "Error while updating tweet", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Tweet updated successfully", http.StatusOK, tweet)
}

// DeleteTweet godoc
// @Router       /api/v1/tweet/{id} [DELETE]
// @Summary      Delete a tweet
// @Description  Delete a tweet by ID
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        id path string true "Tweet ID"
// @Success      204  {object}  nil
// @Failure      403  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) DeleteTweet(c *gin.Context) {
	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	isOwner, err := h.services.Tweet().IsTweetOwner(context.Background(), tweetID, userID)
	if err != nil {
		handleResponse(c, h.log, "Tweet not found or user is not the owner", http.StatusForbidden, err)
		return
	}

	if !isOwner {
		handleResponse(c, h.log, "Unauthorized to update this tweet", http.StatusForbidden, nil)
		return
	}

	err = h.services.Tweet().Delete(context.Background(), tweetID)
	if err != nil {
		handleResponse(c, h.log, "Error while deleting tweet", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Tweet deleted successfully", http.StatusNoContent, nil)
}

// GetTweetList godoc
// @Router       /api/v1/tweets [GET]
// @Summary      Get list of tweets
// @Description  Get list of tweets
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Limit"
// @Success      200  {object}  models.TweetsResponse
// @Failure      500  {object}  models.Response
func (h *Handler) GetTweetList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	request := models.GetListRequest{
		Page:  atoi(page),
		Limit: atoi(limit),
	}

	tweets, err := h.services.Tweet().GetList(context.Background(), request)
	if err != nil {
		handleResponse(c, h.log, "Error while getting tweet list", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Success", http.StatusOK, tweets)
}

// GetTweetsByUser godoc
// @Router       /api/v1/tweets/user/{user_id} [GET]
// @Summary      Get tweets by user
// @Description  Get tweets by user by their User ID
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Param        user_id path string true "User ID"
// @Success      200  {object}  models.TweetsResponse
// @Failure      403  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) GetTweetsByUser(c *gin.Context) {
	userID := c.Param("user_id") // Correct parameter extraction

	fmt.Println(userID) // Debugging print

	// Call service layer to get tweets for the user
	tweets, err := h.services.Tweet().GetTweetsByUser(context.Background(), userID)
	if err != nil {
		handleResponse(c, h.log, "Error while getting tweets by user", http.StatusInternalServerError, err)
		return
	}

	// Send a success response with the tweets
	handleResponse(c, h.log, "Success", http.StatusOK, tweets)
}

// IncrementTweetViews godoc
// @Router       /api/v1/tweet/{id}/views [Patch]
// @Summary      Increment views for a tweet
// @Description  Increment views for a tweet by the logged-in user
// @Tags         tweet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Tweet ID"
// @Success      200  {object}  nil
// @Failure      500  {object}  models.Response
// @Failure      403  {object}  models.Response
// @Failure      401  {object}  models.Response
func (h *Handler) IncrementTweetViews(c *gin.Context) {
	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.services.Tweet().IncrementTweetViews(context.Background(), userID, tweetID)
	if err != nil {
		handleResponse(c, h.log, "Error while incrementing tweet views", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Success", http.StatusOK, nil)
}
