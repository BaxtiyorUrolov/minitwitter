package handler

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"twitter/api/models"
	"twitter/pkg/check"
	"twitter/pkg/email"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var cacheStore = cache.New(5*time.Minute, 10*time.Minute)

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Register godoc
// @Router       /api/v1/register [POST]
// @Summary      Register a new user
// @Description  Register a new user
// @Tags         register
// @Accept       json
// @Produce      json
// @Param        user body models.CreateUser false "user"
// @Success      201  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) Register(c *gin.Context) {
	createUser := models.CreateUser{}
	if err := c.ShouldBindJSON(&createUser); err != nil {
		handleResponse(c, h.log, "Error while reading body from client", http.StatusBadRequest, err)
		return
	}

	if !check.ValidatePassword(createUser.Password) {
		handleResponse(c, h.log, "Invalid password", http.StatusBadRequest, nil)
		return
	}

	exists, err := h.services.User().IsLoginExist(context.Background(), createUser.UserName)
	if err != nil {
		handleResponse(c, h.log, "Error while checking login existence", http.StatusInternalServerError, nil)
		return
	}
	if exists {
		handleResponse(c, h.log, "Login already exists", http.StatusBadRequest, "This login already exists")
		return
	}

	code := generateCode()
	cacheStore.Set(createUser.Email, code, cache.DefaultExpiration)
	cacheStore.Set(createUser.Email+"_data", createUser, cache.DefaultExpiration)
	if err := email.SendEmail(createUser.Email, code); err != nil {
		handleResponse(c, h.log, "Error while sending SMS code", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "SMS kod yuborildi", http.StatusCreated, gin.H{
		"message": "SMS code sent, please verify",
		"phone":   createUser.Email,
	})
}

func (h *Handler) verifyCode(verification models.VerifyCodeRequest) (*models.CreateUser, error) {
	code, found := cacheStore.Get(verification.Email)
	if !found {
		return nil, fmt.Errorf("Verification code expired or not found")
	}
	if code != verification.Code {
		return nil, fmt.Errorf("Invalid verification code")
	}

	// Retrieve cached user data
	userRequest, found := cacheStore.Get(verification.Email + "_data")
	if !found {
		return nil, fmt.Errorf("User data not found")
	}

	createUser := userRequest.(models.CreateUser)
	return &createUser, nil
}

// VerifyRegister godoc
// @Router       /api/v1/verify-register [POST]
// @Summary      Verifies the SMS code
// @Description  verify the SMS code sent to user
// @Tags         register
// @Accept       json
// @Produce      json
// @Param        verification body models.VerifyCodeRequest false "verification"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) VerifyRegister(c *gin.Context) {
	var verification models.VerifyCodeRequest
	if err := c.ShouldBindJSON(&verification); err != nil {
		handleResponse(c, h.log, "Error reading body", http.StatusBadRequest, err)
		return
	}

	createUser, err := h.verifyCode(verification)
	if err != nil {
		handleResponse(c, h.log, "Verification failed", http.StatusBadRequest, err)
		return
	}

	// Create new user after successful verification
	user, err := h.services.User().Create(context.Background(), *createUser)
	if err != nil {
		handleResponse(c, h.log, "Error creating user", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "User created successfully", http.StatusCreated, user)
}
