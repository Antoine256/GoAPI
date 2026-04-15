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
        SELECT id, name, email, role, password, created_at, updated_at, arrival_day, arrival_time, departure_day, departure_time, arrival_info
        FROM users
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []ressources.User
	for rows.Next() {
		var u ressources.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.ArrivalDay, &u.ArrivalTime, &u.DepartureDay, &u.DepartureTime, &u.ArrivalInfo); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserByID(id int) (ressources.User, error) {
	var u ressources.User
	err := database.DB.QueryRow(`
        SELECT id, name, email, role, password, created_at, updated_at, arrival_day, arrival_time, departure_day, departure_time, arrival_info
        FROM users WHERE id = $1
    `, id).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.ArrivalDay, &u.ArrivalTime, &u.DepartureDay, &u.DepartureTime, &u.ArrivalInfo)
	if errors.Is(err, sql.ErrNoRows) {
		return ressources.User{}, errors.New("utilisateur introuvable")
	}
	return u, err
}

func GetUserByEmail(email string) (ressources.User, error) {
	var u ressources.User
	err := database.DB.QueryRow(`
        SELECT id, name, email, role, password, created_at, updated_at, arrival_day, arrival_time, departure_day, departure_time, arrival_info
        FROM users WHERE email = $1
    `, email).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.ArrivalDay, &u.ArrivalTime, &u.DepartureDay, &u.DepartureTime, &u.ArrivalInfo)
	if errors.Is(err, sql.ErrNoRows) {
		return ressources.User{}, errors.New("utilisateur introuvable")
	}
	return u, err
}

func CreateUser(u ressources.User) (ressources.User, error) {
	// Scan permet de récupérer les champs auto-générés (id, timestamps) après l'insertion
	err := database.DB.QueryRow(`
        INSERT INTO users (name, email, password, role)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at, arrival_day, arrival_time, departure_day, departure_time, arrival_info
    `, u.Name, u.Email, u.Password, u.Role).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.ArrivalDay, &u.ArrivalTime, &u.DepartureDay, &u.DepartureTime, &u.ArrivalInfo)
	return u, err
}

func UpdateUser(u ressources.User) (ressources.User, error) {
	err := database.DB.QueryRow(`
        UPDATE users
        SET name = $1, email = $2, password = $3, role = $4, updated_at = $5, arrival_day = $6, arrival_time = $7, departure_day = $8, departure_time = $9, arrival_info = $10
        WHERE id = $11
        RETURNING updated_at
    `, u.Name, u.Email, u.Password, u.Role, time.Now(), u.ArrivalDay, u.ArrivalTime, u.DepartureDay, u.DepartureTime, u.ArrivalInfo, u.ID).Scan(&u.UpdatedAt)
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
