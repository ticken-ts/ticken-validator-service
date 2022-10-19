package services

import (
	"ticken-validator-service/models"
)

type IProvider interface {
	GetTicketScanner() TicketScanner
	GetEventManager() EventManager
}

type TicketScanner interface {
	Scan(eventID string, ticketID string, owner string) (*models.Ticket, error)
}

type EventManager interface {
	AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error)
}
