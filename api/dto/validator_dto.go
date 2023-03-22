package dto

type ValidatorDTO struct {
	ValidatorID    string `json:"validator_id"`
	OrganizationID string `json:"organization_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
}
