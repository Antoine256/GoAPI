package repository

import (
	"GoAPI/database"
	"database/sql"
	"errors"
)

func SaveRefreshToken(userID int, token string) error {
	_, err := database.DB.Exec(`
        INSERT INTO refresh_tokens (user_id, token)
        VALUES ($1, $2)
    `, userID, token)
	return err
}

func RefreshTokenExists(userID int, token string) (bool, error) {
	var id int
	err := database.DB.QueryRow(`
        SELECT id FROM refresh_tokens
        WHERE user_id = $1 AND token = $2
    `, userID, token).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteRefreshToken(token string) error {
	_, err := database.DB.Exec(`
        DELETE FROM refresh_tokens WHERE token = $1
    `, token)
	return err
}
