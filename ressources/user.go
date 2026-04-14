package ressources

type User struct {
	ID            int
	Name          string
	Email         string
	Password      string
	Role          string
	CreatedAt     string
	UpdatedAt     string
	ArrivalDay    string
	ArrivalTime   string
	DepartureDay  string
	DepartureTime string
}

type UserPublicDTO struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ArrivalDay    string `json:"arrival_day"`
	ArrivalTime   string `json:"arrival_time"`
	DepartureDay  string `json:"departure_day"`
	DepartureTime string `json:"departure_time"`
}

func ToUserPublicDTO(user User) UserPublicDTO {
	return UserPublicDTO{
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
	}
}

func ToUser(request UserPublicDTO) User {
	return User{
		ID:            request.ID,
		Name:          request.Name,
		Email:         request.Email,
		Role:          request.Role,
		CreatedAt:     request.CreatedAt,
		UpdatedAt:     request.UpdatedAt,
		ArrivalDay:    request.ArrivalDay,
		ArrivalTime:   request.ArrivalTime,
		DepartureDay:  request.DepartureDay,
		DepartureTime: request.DepartureTime,
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
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
	Role          string `json:"role" binding:"required"`
	ArrivalDay    string `json:"arrival_day"`
	ArrivalTime   string `json:"arrival_time"`
	DepartureDay  string `json:"departure_day"`
	DepartureTime string `json:"departure_time"`
}

func (dto UserUpdateDTO) ToUser() User {
	return User{
		Name:          dto.Name,
		Email:         dto.Email,
		Password:      dto.Password,
		Role:          dto.Role,
		ArrivalDay:    dto.ArrivalDay,
		ArrivalTime:   dto.ArrivalTime,
		DepartureDay:  dto.DepartureDay,
		DepartureTime: dto.DepartureTime,
	}
}
