package services

import (
	"github.com/google/uuid"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
)

type AttendantManager struct {
	attendantRepo repos.IAttendantRepository
}

func NewAttendantManager(repoProvider repos.IProvider) *AttendantManager {
	return &AttendantManager{attendantRepo: repoProvider.GetAttendantRepository()}
}

func (attendantManager *AttendantManager) AddAttendant(attendantID uuid.UUID, walletAddress string, pemPublicKey string) (*models.Attendant, error) {
	attendant := &models.Attendant{
		AttendantID:   attendantID,
		WalletAddress: walletAddress,
		PemPublicKey:  pemPublicKey,
	}

	if err := attendantManager.attendantRepo.AddAttendant(attendant); err != nil {
		return nil, err
	}
	return attendant, nil
}
