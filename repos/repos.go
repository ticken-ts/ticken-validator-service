package repos

import (
	"github.com/google/uuid"
	"ticken-validator-service/models"
)

type IEventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID uuid.UUID) *models.Event
}

type ITicketRepository interface {
	AddTicket(ticket *models.Ticket) error
	AddManyTickets(ticket []*models.Ticket) error
	FindTicket(eventID uuid.UUID, ticketID uuid.UUID) *models.Ticket
}

type IAttendantRepository interface {
	AddAttendant(attendant *models.Attendant) error
	FindAttendant(attendantID uuid.UUID) *models.Attendant
	FindAttendantByWalletAddr(wallerAddr string) *models.Attendant
}

type IProvider interface {
	GetEventRepository() IEventRepository
	GetTicketRepository() ITicketRepository
	GetAttendantRepository() IAttendantRepository
}

type IFactory interface {
	BuildEventRepository() any
	BuildTicketRepository() any
	BuildAttendantRepository() any
}
