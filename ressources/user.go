package ressources

import "time"

// Modèle DB
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
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
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
}

func (dto UserUpdateRequest) ToUser() User {
	return User{
		Name:  *dto.Name,
		Email: *dto.Email,
		Role:  *dto.Role,
	}
}

// Response

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (user User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
