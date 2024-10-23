package handler

import (
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services service.IServiceManager
	log      logger.Logger
}

func New(services service.IServiceManager, log logger.Logger) Handler {
	return Handler{
		services: services,
		log:      log,
	}
}

func handleResponse(c *gin.Context, log logger.Logger, msg string, statusCode int, data interface{}) {
	resp := models.Response{}

	switch code := statusCode; {
	case code < 400:
		resp.Description = "OK"
		log.Info("~~~~> OK", logger.String("msg", msg), logger.Any("status", code))
	case code < 500:
		resp.Description = "Bad Request"
		log.Error("!!!!! BAD REQUEST", logger.String("msg", msg), logger.Any("status", code))
	default:
		resp.Description = "Internal Server Error"
		log.Error("!!!!! INTERNAL SERVER ERROR", logger.String("msg", msg), logger.Any("status", code))
	}

	resp.StatusCode = statusCode
	resp.Data = data

	c.JSON(resp.StatusCode, resp)
}
