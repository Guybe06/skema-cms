package types

type ProfileResponse struct {
	ID            string `json:"id"             example:"018f1e2a-3b4c-7d8e-9f0a-1b2c3d4e5f6a"`
	Email         string `json:"email"          example:"jane@skemacms.com"`
	FirstName     string `json:"first_name"     example:"Jane"`
	LastName      string `json:"last_name"      example:"Doe"`
	EmailVerified bool   `json:"email_verified" example:"true"`
	CreatedAt     string `json:"created_at"     example:"2026-01-15T10:30:00Z"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=1" example:"Jane"`
	LastName  string `json:"last_name"  validate:"omitempty,min=1" example:"Doe"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"       example:"OldP@ss123!"`
	NewPassword     string `json:"new_password"     validate:"required,min=8" example:"N3wP@ss456!"`
}
