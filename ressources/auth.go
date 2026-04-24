package ressources

// DTOs

// Requests

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SecretKeyRequest struct {
	Key string `json:"key" binding:"required"`
}

// Responses

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
