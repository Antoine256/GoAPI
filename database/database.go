package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var DB *sql.DB

func Connect(logger *zap.Logger) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal("Impossible d'ouvrir la connexion : %v", zap.Error(err))
	}

	if err = DB.Ping(); err != nil {
		logger.Fatal("Impossible de joindre la base : %v", zap.Error(err))
	}

	logger.Info("Connecté à PostgreSQL")
	migrate(logger)
}

func migrate(logger *zap.Logger) {
	// Implémentation de la logique de migration (ex: création de tables)
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id         SERIAL PRIMARY KEY,
            name       VARCHAR(100) NOT NULL UNIQUE,
            email      VARCHAR(150) NOT NULL UNIQUE,
			role       VARCHAR(20) NOT NULL DEFAULT 'user',
            password   VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
			arrival_info BOOLEAN NOT NULL DEFAULT FALSE,
			arrival_day VARCHAR(20) NOT NULL DEFAULT 'non renseigné',
			arrival_time VARCHAR(20) NOT NULL DEFAULT 'non renseigné',
			departure_day VARCHAR(20) NOT NULL DEFAULT 'non renseigné',
			departure_time VARCHAR(20) NOT NULL DEFAULT 'non renseigné'
        )`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
            id         SERIAL PRIMARY KEY,
            user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            token      TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        )`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			logger.Fatal("Erreur migration : %v", zap.Error(err))
		}
	}

	logger.Info("Migration effectuée")

}
