package mappers

import (
	"ticken-validator-service/api/dto"
	"ticken-validator-service/models"
)

func MapValidatorToDTO(validator *models.Validator) *dto.ValidatorDTO {
	return &dto.ValidatorDTO{
		ValidatorID:    validator.ValidatorID.String(),
		OrganizationID: validator.OrganizationID.String(),
		Username:       validator.Username,
		Email:          validator.Email,
	}
}
