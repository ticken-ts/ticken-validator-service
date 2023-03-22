package services

import (
	"fmt"
	"github.com/google/uuid"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
	"ticken-validator-service/security/jwt"
)

type ValidatorManager struct {
	jwtVerifier   jwt.Verifier
	validatorRepo repos.IValidatorRepository
}

func NewValidatorManager(repoProvider repos.IProvider, jwtVerifier jwt.Verifier) *ValidatorManager {
	return &ValidatorManager{jwtVerifier: jwtVerifier, validatorRepo: repoProvider.GetValidatorRepository()}
}

func (manager *ValidatorManager) RegisterValidator(organizationID uuid.UUID, validatorJWT string) (*models.Validator, error) {
	jwtToken, err := manager.jwtVerifier.Verify(validatorJWT)
	if err != nil {
		return nil, fmt.Errorf("failed to verify validators JWT: %s", err)
	}

	newValidator := &models.Validator{
		ValidatorID:    jwtToken.Subject,
		Username:       jwtToken.Username,
		Email:          jwtToken.Email,
		OrganizationID: organizationID,
	}

	validatorWithSameID := manager.validatorRepo.FindValidator(newValidator.ValidatorID)
	if validatorWithSameID != nil {
		return nil, fmt.Errorf("validator with ID %s already exists", validatorWithSameID.ValidatorID)
	}

	if err := manager.validatorRepo.AddValidator(newValidator); err != nil {
		return nil, err
	}

	return newValidator, nil
}
