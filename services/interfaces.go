package services

import (
	"github.com/google/uuid"
	"ticken-validator-service/models"
)

type IProvider interface {
	GetTicketScanner() ITicketScanner
	GetTicketSyncer() ITicketSyncer
	GetEventManager() IEventManager
	GetAttendantManager() IAttendantManager
	GetValidatorManager() IValidatorManager
}

type ITicketScanner interface {
	Scan(eventID, ticketID uuid.UUID, signature []byte, validatorID uuid.UUID) (*models.Ticket, error)
}

type ITicketSyncer interface {
	Sync(eventID uuid.UUID, callerID uuid.UUID) error
}

type IEventManager interface {
	AddEvent(eventID, organizerID, organizationID uuid.UUID, pvtBCChannel, pubBCAddress string) (*models.Event, error)
}

type IAttendantManager interface {
	AddAttendant(attendantID uuid.UUID, walletAddress string, publicKey []byte) (*models.Attendant, error)
}

type IValidatorManager interface {
	RegisterValidator(organizationID uuid.UUID, validatorJWT string) (*models.Validator, error)
}
