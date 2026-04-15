package ressources

// DTOs

// Requests

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (dto LoginRequest) ToUser() User {
	return User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (dto RegisterRequest) ToUser() User {
	return User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

// Responses

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
