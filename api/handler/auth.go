package handler

import (
	"context"
	"errors"
	"net/http"
	"time"
	"twitter/api/models"
	"twitter/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Router       /api/v1/login [POST]
// @Summary      Customer login
// @Description  Customer login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.LoginRequest false "login"
// @Success      201  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) Login(c *gin.Context) {
	loginRequest := models.LoginRequest{}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		handleResponse(c, h.log, "error while binding body", http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	loginResponse, err := h.services.Auth().Login(ctx, loginRequest)
	if err != nil {
		handleResponse(c, h.log, "incorrect credentials", http.StatusBadRequest, errors.New("password or login incorrect"))
		return
	}

	handleResponse(c, h.log, "success", http.StatusOK, loginResponse)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		m, err := jwt.ExtractClaims(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", m["user_id"].(string))
		c.Next()
	}
}
