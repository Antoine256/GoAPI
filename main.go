package main

import (
	"GoAPI/database"
	"GoAPI/router"
	logger "GoAPI/utils"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Impossible de charger le fichier .env")
		if os.Getenv("DB_HOST") != "" {
			fmt.Println("Variables d'environnement trouvées dans le système")
		}
	}
	logger, err := logger.NewLogger("info")
	if err != nil {
		log.Fatalf("Erreur initialisation logger : %v", err)
	}
	defer logger.Sync()

	database.Connect(logger)

	// Create a Gin router with default middleware (logger and recovery)
	r := router.SetupRouter(logger)

	// Start server on port 8690
	if err := r.Run("0.0.0.0:8690"); err != nil {
		logger.Fatal("failed to run server: %v", zap.Error(err))
	}
}
