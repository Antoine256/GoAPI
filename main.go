package main

import (
	"GoAPI/database"
	"GoAPI/router"
	logger "GoAPI/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func helloworld(c *gin.Context) {
	c.String(http.StatusOK, "Hello World! Time : %s", time.Now().Format(time.RFC3339Nano))
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur chargement .env")
	}

	database.Connect()

	logger, err := logger.NewLogger("info")
	if err != nil {
		log.Fatalf("Erreur initialisation logger : %v", err)
	}
	defer logger.Sync()

	// Create a Gin router with default middleware (logger and recovery)
	r := router.SetupRouter(logger)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
