package services

import (
	"GoAPI/repository"
	"GoAPI/ressources"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers() ([]ressources.User, error) {
	return repository.GetAllUsers()
}

func GetUserByID(id int, logger *zap.Logger) (ressources.User, error) {
	return repository.GetUserByID(id)
}

func CreateUser(dto ressources.UserCreateRequest, logger *zap.Logger) (ressources.User, error) {
	// Hash du mot de passe avant insertion
	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("CreateUser - failed to hash password", zap.Error(err))
		return ressources.User{}, err
	}

	user := dto.ToUser()
	user.Password = string(hashed)

	return repository.CreateUser(user)
}

func UpdateUser(id int, dto ressources.UserUpdateRequest, logger *zap.Logger) (ressources.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("UpdateUser - failed to hash password", zap.Error(err))
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
