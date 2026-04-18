package ressources

import "time"

// Modèle DB
type User struct {
	ID            int
	Name          string
	Email         string
	Password      string
	Role          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ArrivalInfo   bool
	ArrivalDay    string
	ArrivalTime   string
	DepartureDay  string
	DepartureTime string
}

// DTOs

// Request

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (dto UserCreateRequest) ToUser() User {
	return User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Role:     dto.Role,
	}
}

type UserUpdateRequest struct {
	Name          *string `json:"name"`
	Email         *string `json:"email"`
	Role          *string `json:"role"`
	ArrivalInfo   *bool   `json:"arrival_info"`
	ArrivalDay    *string `json:"arrival_day"`
	ArrivalTime   *string `json:"arrival_time"`
	DepartureDay  *string `json:"departure_day"`
	DepartureTime *string `json:"departure_time"`
}

func (dto UserUpdateRequest) ToUser() User {
	return User{
		Name:          *dto.Name,
		Email:         *dto.Email,
		Role:          *dto.Role,
		ArrivalInfo:   *dto.ArrivalInfo,
		ArrivalDay:    *dto.ArrivalDay,
		ArrivalTime:   *dto.ArrivalTime,
		DepartureDay:  *dto.DepartureDay,
		DepartureTime: *dto.DepartureTime,
	}
}

// Response

type UserResponse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ArrivalDay    string    `json:"arrival_day"`
	ArrivalTime   string    `json:"arrival_time"`
	DepartureDay  string    `json:"departure_day"`
	DepartureTime string    `json:"departure_time"`
	ArrivalInfo   bool      `json:"arrival_info"`
}

func (user User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		ArrivalDay:    user.ArrivalDay,
		ArrivalTime:   user.ArrivalTime,
		DepartureDay:  user.DepartureDay,
		DepartureTime: user.DepartureTime,
		ArrivalInfo:   user.ArrivalInfo,
	}
}
