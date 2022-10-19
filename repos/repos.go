package repos

import "ticken-validator-service/models"

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID string) *models.Event
}

type TicketRepository interface {
	AddTicket(ticket *models.Ticket) error
	FindTicket(eventID string, ticketID string) *models.Ticket
}

type IProvider interface {
	GetEventRepository() EventRepository
	GetTicketRepository() TicketRepository
}

type IFactory interface {
	BuildEventRepository() any
	BuildTicketRepository() any
}
