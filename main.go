package main

import (
	"GoAPI/database"
	"GoAPI/router"
	logger "GoAPI/utils"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

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

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Start server on port 8690
	if err := r.Run("localhost:8690"); err != nil {
		logger.Fatal("failed to run server: %v", zap.Error(err))
	}
}
