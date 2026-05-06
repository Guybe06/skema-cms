package types

type InviteRequest struct {
	Email string `json:"email" validate:"required,email" example:"colleague@skemacms.com"`
	Role  string `json:"role"  validate:"required,oneof=admin member" example:"member"`
}

type AcceptInviteRequest struct {
	Token string `json:"token" validate:"required" example:"a3f2c1d4e5b6..."`
}

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin member" example:"admin"`
}

type MemberResponse struct {
	ID        string `json:"id"         example:"018f1e2a-3b4c-7d8e-9f0a-1b2c3d4e5f6a"`
	UserID    string `json:"user_id"    example:"018f1e2a-3b4c-7d8e-9f0a-1b2c3d4e5f6b"`
	Email     string `json:"email"      example:"colleague@skemacms.com"`
	Role      string `json:"role"       example:"member"`
	Status    string `json:"status"     example:"active"`
	InvitedBy string `json:"invited_by" example:"018f1e2a-3b4c-7d8e-9f0a-1b2c3d4e5f6c"`
	CreatedAt string `json:"created_at" example:"2026-01-15T10:30:00Z"`
}
