package services

import (
	"GoAPI/repository"
	"GoAPI/ressources"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers() ([]ressources.UserResponse, error) {
	// on convertit les users en UserResponse pour ne pas exposer les mots de passe
	users, err := repository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []ressources.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToUserResponse())
	}

	return userResponses, nil
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

	oldUser, err := GetUserByID(id, logger)
	if err != nil {
		logger.Error("UpdateUser - user not found", zap.Int("id", id), zap.Error(err))
		return ressources.User{}, err
	}

	if dto.Name != nil {
		oldUser.Name = *dto.Name
	}
	if dto.Email != nil {
		oldUser.Email = *dto.Email
	}
	if dto.Role != nil {
		oldUser.Role = *dto.Role
	}
	if dto.ArrivalInfo != nil {
		oldUser.ArrivalInfo = *dto.ArrivalInfo
	}
	if dto.ArrivalDay != nil {
		oldUser.ArrivalDay = *dto.ArrivalDay
	}
	if dto.ArrivalTime != nil {
		oldUser.ArrivalTime = *dto.ArrivalTime
	}
	if dto.DepartureDay != nil {
		oldUser.DepartureDay = *dto.DepartureDay
	}
	if dto.DepartureTime != nil {
		oldUser.DepartureTime = *dto.DepartureTime
	}

	if dto.ArrivalDay != nil && dto.ArrivalTime != nil && dto.DepartureDay != nil && dto.DepartureTime != nil {
		oldUser.ArrivalInfo = true
	}

	logger.Info("UpdateUser - updating user", zap.Int("id", id), zap.Any("updated_fields", dto))

	return repository.UpdateUser(oldUser)
}

func DeleteUser(id int) error {
	return repository.DeleteUser(id)
}
