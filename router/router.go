package router

import (
	"GoAPI/handlers"
	middleware "GoAPI/middlewares"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	logger.Info("Allowed origins for CORS", zap.Strings("origins", allowedOrigins))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(gin.Recovery())

	// Routes publiques

	authHandler := handlers.NewAuthHandler(logger)

	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}

	// Routes protégées

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(logger))

	userHandler := handlers.NewUserHandler(logger)

	users := api.Group("/users")
	{
		users.GET("", userHandler.GetUsers)
		users.GET("/me", userHandler.GetCurrentUser)
		users.GET("/:id", userHandler.GetUser)
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	return r
}
