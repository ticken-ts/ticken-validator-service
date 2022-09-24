package services

import "ticken-validator-service/models"

type Provider interface {
	GetTicketScanner() TicketScanner
}

type TicketScanner interface {
	Scan(eventID string, ticketID string, owner string) (*models.Ticket, error)
}
