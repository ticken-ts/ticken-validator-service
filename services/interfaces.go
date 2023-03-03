package services

import (
	"github.com/google/uuid"
	"ticken-validator-service/models"
)

type IProvider interface {
	GetTicketScanner() ITicketScanner
	GetTicketSyncer() ITicketSyncer
	GetEventManager() IEventManager
}

type ITicketScanner interface {
	Scan(eventID, ticketID uuid.UUID, signature string, validatorID uuid.UUID) (*models.Ticket, error)
}

type ITicketSyncer interface {
	Sync(eventID uuid.UUID, callerID uuid.UUID) error
}

type IEventManager interface {
	AddEvent(eventID, organizerID uuid.UUID, pvtBCChannel, pubBCAddress string) (*models.Event, error)
}
