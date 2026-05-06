package types

type RegisterRequest struct {
	Email     string `json:"email"      validate:"required,email"  example:"jane@skemacms.com"`
	Password  string `json:"password"   validate:"required,min=8"  example:"S3cur3P@ss!"`
	FirstName string `json:"first_name" validate:"required"         example:"Jane"`
	LastName  string `json:"last_name"  validate:"required"         example:"Doe"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"jane@skemacms.com"`
	Password string `json:"password" validate:"required"       example:"S3cur3P@ss!"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required" example:"a3f2c1d4e5b6..."`
}

type RequestResetRequest struct {
	Email string `json:"email" validate:"required,email" example:"jane@skemacms.com"`
}

type ConfirmResetRequest struct {
	Token    string `json:"token"    validate:"required"        example:"a3f2c1d4e5b6..."`
	Password string `json:"password" validate:"required,min=8"  example:"N3wS3cur3P@ss!"`
}

type UserResponse struct {
	ID            string `json:"id"             example:"018f1e2a-3b4c-7d8e-9f0a-1b2c3d4e5f6a"`
	Email         string `json:"email"          example:"jane@skemacms.com"`
	FirstName     string `json:"first_name"     example:"Jane"`
	LastName      string `json:"last_name"      example:"Doe"`
	EmailVerified bool   `json:"email_verified" example:"false"`
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int          `json:"expires_in"    example:"3600"`
	User         UserResponse `json:"user"`
}
