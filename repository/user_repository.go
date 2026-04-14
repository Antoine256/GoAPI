package repository

import (
	"GoAPI/database"
	"GoAPI/ressources"
	"database/sql"
	"errors"
	"time"
)

func GetAllUsers() ([]ressources.User, error) {
	rows, err := database.DB.Query(`
        SELECT id, name, email, role, password, created_at, updated_at
        FROM users
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []ressources.User
	for rows.Next() {
		var u ressources.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserByID(id int) (ressources.User, error) {
	var u ressources.User
	err := database.DB.QueryRow(`
        SELECT id, name, email, role, password, created_at, updated_at
        FROM users WHERE id = $1
    `, id).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ressources.User{}, errors.New("utilisateur introuvable")
	}
	return u, err
}

func GetUserByEmail(email string) (ressources.User, error) {
	var u ressources.User
	err := database.DB.QueryRow(`
        SELECT id, name, email, role, password, created_at, updated_at
        FROM users WHERE email = $1
    `, email).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ressources.User{}, errors.New("utilisateur introuvable")
	}
	return u, err
}

func CreateUser(u ressources.User) (ressources.User, error) {
	err := database.DB.QueryRow(`
        INSERT INTO users (name, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `, u.Name, u.Email, u.Password).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func UpdateUser(u ressources.User) (ressources.User, error) {
	err := database.DB.QueryRow(`
        UPDATE users
        SET name = $1, email = $2, password = $3, role = $4, updated_at = $5
        WHERE id = $6
        RETURNING updated_at
    `, u.Name, u.Email, u.Password, u.Role, time.Now(), u.ID).Scan(&u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ressources.User{}, errors.New("utilisateur introuvable")
	}
	return u, err
}

func DeleteUser(id int) error {
	result, err := database.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("utilisateur introuvable")
	}
	return nil
}
