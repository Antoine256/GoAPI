package services

import (
	"GoAPI/repository"
	"GoAPI/ressources"

	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers() ([]ressources.User, error) {
	return repository.GetAllUsers()
}

func GetUserByID(id int) (ressources.User, error) {
	return repository.GetUserByID(id)
}

func CreateUser(dto ressources.UserCreateDTO) (ressources.User, error) {
	// Hash du mot de passe avant insertion
	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return ressources.User{}, err
	}

	user := dto.ToUser()
	user.Password = string(hashed)

	return repository.CreateUser(user)
}

func UpdateUser(id int, dto ressources.UserUpdateDTO) (ressources.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return ressources.User{}, err
	}

	user := dto.ToUser()
	user.ID = id
	user.Password = string(hashed)

	return repository.UpdateUser(user)
}

func DeleteUser(id int) error {
	return repository.DeleteUser(id)
}
