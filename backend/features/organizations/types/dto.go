package types

type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255" example:"Acme Corp"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=255" example:"Acme Corp Updated"`
}

type TransferOwnershipRequest struct {
	NewOwnerID string `json:"new_owner_id" validate:"required,uuid4" example:"018f1e2a-3b4c-7d8e-9f0a-1b2c3d4e5f6b"`
}

type OrganizationResponse struct {
	ID        string `json:"id"         example:"018f1e2a-3b4c-7d8e-9f0a-1b2c3d4e5f6a"`
	Name      string `json:"name"       example:"Acme Corp"`
	Slug      string `json:"slug"       example:"acme-corp"`
	OwnerID   string `json:"owner_id"   example:"018f1e2a-3b4c-7d8e-9f0a-1b2c3d4e5f6b"`
	CreatedAt string `json:"created_at" example:"2026-01-15T10:30:00Z"`
}
