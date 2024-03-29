package async

import (
	"encoding/json"
	"github.com/google/uuid"
	"ticken-validator-service/log"
	"ticken-validator-service/services"
)

const (
	NewAttendantMessageType = "new_attendant"
)

type attendantDTO struct {
	AttendantID   uuid.UUID `json:"attendant_id"`
	WalletAddress string    `json:"wallet_address"`
	PublicKey     string    `json:"public_key"`
}

type AttendantSubscriber struct {
	attendantManager services.IAttendantManager
}

func NewAttendantSubscriber(attendantManager services.IAttendantManager) *AttendantSubscriber {
	return &AttendantSubscriber{attendantManager: attendantManager}
}

func (s *AttendantSubscriber) NewAttendantHandler(rawAttendant []byte) error {
	dto := new(attendantDTO)

	log.TickenLogger.Info().Msg("loading attendant")

	err := json.Unmarshal(rawAttendant, dto)
	if err != nil {
		return err
	}

	_, err = s.attendantManager.AddAttendant(dto.AttendantID, dto.WalletAddress, dto.PublicKey)
	if err != nil {
		return err
	}

	return nil
}
