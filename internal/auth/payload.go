package auth

type LoginRequestDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=5"`
}
type AuthRequestDTO struct {
	Username string `json:"username" validate:"required,min=3,max=25"`
	Password string `json:"password" validate:"required,min=5"`
}

type AuthResponseDTO struct {
	AccessToken string `json:"access_token"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=25"`
	Password string `json:"password" validate:"required,min=5"`
}

type RegisterResponse struct {
	RegisterSuccess bool   `json:"registerSuccess"`
	Message         string `json:"message"`
}
