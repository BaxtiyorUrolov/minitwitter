package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"twitter/api/models"
	"twitter/pkg/check"
)

// CreateUser godoc
// @Router       /api/v1/user [POST]
// @Summary      Creates a new user
// @Description  create a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user body models.CreateUser false "user"
// @Success      201  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CreateUser(c *gin.Context) {
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

}

// GetUser godoc
// @Router       /api/v1/user/{id} [GET]
// @Summary      Get a user by ID
// @Description  Get a user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  models.User
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.services.User().GetByID(context.Background(), models.PrimaryKey{ID: id})
	if err != nil {
		handleResponse(c, h.log, "User not found", http.StatusNotFound, err)
		return
	}

	handleResponse(c, h.log, "Success", http.StatusOK, user)
}

// GetUserList godoc
// @Router       /api/v1/users [GET]
// @Summary      Get list of users
// @Description  Get list of users
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Limit"
// @Success      200  {object}  models.UsersResponse
// @Failure      500  {object}  models.Response
func (h *Handler) GetUserList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	request := models.GetListRequest{
		Page:  atoi(page),
		Limit: atoi(limit),
	}

	users, err := h.services.User().GetAllUsers(context.Background(), request)
	if err != nil {
		handleResponse(c, h.log, "Error while getting user list", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "Success", http.StatusOK, users)
}

// UpdateUser godoc
// @Router       /api/v1/user/{id} [PUT]
// @Summary      Update a user
// @Description  Update a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        user body models.UpdateUser false "user"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) UpdateUser(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		handleResponse(c, h.log, "Unauthorized", http.StatusForbidden, "User not authenticated")
		return
	}

	updateUser := models.UpdateUser{ID: userID.(string)}
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		handleResponse(c, h.log, "Error while binding request body", http.StatusBadRequest, err)
		return
	}

	err := h.services.User().Update(context.Background(), updateUser)
	if err != nil {
		handleResponse(c, h.log, "Error while updating user", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "User updated successfully", http.StatusOK, nil)
}

// DeleteUser godoc
// @Router       /api/v1/user/{id} [DELETE]
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiAuthKey
// @Param        id path string true "User ID"
// @Success      204  {object}  nil
// @Failure      403  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handleResponse(c, h.log, "Unauthorized", http.StatusForbidden, "User not authenticated")
		return
	}

	err := h.services.User().Delete(context.Background(), models.PrimaryKey{ID: userID.(string)})
	if err != nil {
		handleResponse(c, h.log, "Error while deleting user", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "User deleted successfully", http.StatusNoContent, nil)
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Error converting string to int: %v", err)
		return 1
	}
	return i
}
