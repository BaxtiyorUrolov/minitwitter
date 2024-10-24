package api

import (
	_ "twitter/api/docs"

	_ "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"twitter/api/handler"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description API key needed to access the endpoints

func (s *Server) endpoints() {
	v1 := s.router.Group("/api/v1") // Group creation
	v1.POST("/register", s.handler.Register)
	v1.POST("/verify-register", s.handler.VerifyRegister)

	v1.POST("/login", s.handler.Login)

	v1.POST("/user", s.handler.CreateUser)
	v1.GET("/user/:id", s.handler.GetUser)
	v1.GET("/users", s.handler.GetUserList)

	// Tweet

	v1.GET("/tweet/:id", s.handler.GetTweet)
	v1.GET("/tweets", s.handler.GetTweetList)
	v1.GET("/tweets/user/:user_id", s.handler.GetTweetsByUser)

	// JWT Middleware protected routes
	protectedRoutes := v1.Group("/") // "/api/v1/" under new group
	protectedRoutes.Use(handler.AuthMiddleware())
	{
		protectedRoutes.PUT("/user/:id", s.handler.UpdateUser)
		protectedRoutes.DELETE("/user/:id", s.handler.DeleteUser)

		// Tweet

		protectedRoutes.POST("/tweet", s.handler.CreateTweet)
		protectedRoutes.PUT("/tweet/:id", s.handler.UpdateTweet)
		protectedRoutes.DELETE("/tweet/:id", s.handler.DeleteTweet)
		protectedRoutes.PATCH("/tweet/:id/views", s.handler.IncrementTweetViews)
	}

	// Swagger documentation
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
