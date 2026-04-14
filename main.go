package main

import (
	"GoAPI/database"
	"GoAPI/router"
	logger "GoAPI/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func helloworld(c *gin.Context) {
	c.String(http.StatusOK, "Hello World! Time : %s", time.Now().Format(time.RFC3339Nano))
}

func main() {

	logger, err := logger.NewLogger("info")
	if err != nil {
		log.Fatalf("Erreur initialisation logger : %v", err)
	}
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Warn("Impossible de charger le fichier .env")
		if os.Getenv("DB_HOST") != "" {
			logger.Info("Variables d'environnement trouvées dans le système")
		}
	}

	database.Connect(logger)

	// Create a Gin router with default middleware (logger and recovery)
	r := router.SetupRouter(logger)

	// Start server on port 8690
	if err := r.Run("localhost:8690"); err != nil {
		logger.Fatal("failed to run server: %v", zap.Error(err))
	}
}
