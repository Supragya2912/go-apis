package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdatePasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type UpdatePasswordResponse struct {
	Message string `json:"message"`
}
