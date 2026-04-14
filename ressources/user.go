package ressources

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt string
	UpdatedAt string
}

type UserPublicDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToUserPublicDTO(user User) UserPublicDTO {
	return UserPublicDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUser(request UserPublicDTO) User {
	return User{
		ID:        request.ID,
		Name:      request.Name,
		Email:     request.Email,
		Role:      request.Role,
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
	}
}

type UserCreateDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (dto UserCreateDTO) ToUser() User {
	return User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type UserUpdateDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (dto UserUpdateDTO) ToUser() User {
	return User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Role:     dto.Role,
	}
}
