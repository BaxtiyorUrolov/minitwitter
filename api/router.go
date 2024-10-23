package api

import (
	_ "twitter/api/docs"

	_ "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"twitter/api/handler"
)

// @securityDefinitions.apikey ApiAuthKey
// @in header
// @name X-API-Key
// @description API key needed to access the endpoints

func (s *Server) endpoints() {
	v1 := s.router.Group("/api/v1") // Group yaratish
	v1.POST("/register", s.handler.Register)
	v1.POST("/verify-register", s.handler.VerifyRegister)

	v1.POST("/login", s.handler.Login)

	v1.POST("/user", s.handler.CreateUser)
	v1.GET("/user/:id", s.handler.GetUser)
	v1.GET("/users", s.handler.GetUserList)
	v1.DELETE("/user/:id", s.handler.DeleteUser)

	// JWT Middleware bilan himoyalangan marshrutlar
	protectedRoutes := v1.Group("/") // "/api/v1/" ostida yangi guruh
	protectedRoutes.Use(handler.AuthMiddleware())
	{
		protectedRoutes.PUT("/user", s.handler.UpdateUser)
		protectedRoutes.DELETE("/user", s.handler.DeleteUser)
	}

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
