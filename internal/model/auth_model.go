package model

type Auth struct {
	ID uint `json:"id"`
}

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthLoginResponse struct {
	Token string `json:"token" validate:"required"`
}

type AuthLogoutRequest struct {
	ID uint `json:"id" validate:"required"`
}

type AuthSignUpRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthVerifyRequest struct {
	Token string `validate:"required"`
}
