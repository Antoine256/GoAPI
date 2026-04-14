package router

import (
	"GoAPI/handlers"
	middleware "GoAPI/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(gin.Recovery())

	// Routes publiques

	auth := r.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
		auth.POST("/refresh", handlers.RefreshToken)
		auth.POST("/logout", handlers.Logout)
	}

	// Routes protégées
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(logger))
	users := api.Group("/users")
	{
		users.GET("", handlers.GetUsers)
		users.GET("/me", handlers.GetCurrentUser)
		users.GET("/:id", handlers.GetUser)
		users.POST("", handlers.CreateUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}

	return r
}
