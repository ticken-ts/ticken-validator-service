package services

import (
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-validator-service/models"
)

type Provider interface {
	GetTicketScanner() TicketScanner
	GetEventManager() EventManager
}

type TicketScanner interface {
	Scan(eventID string, ticketID string, owner string) (*chain_models.Ticket, error)
}

type EventManager interface {
	AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error)
}
