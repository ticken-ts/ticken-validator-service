package mappers

import (
	"ticken-validator-service/api/dto"
	"ticken-validator-service/models"
)

func MapTicketToDTO(ticket *models.Ticket) *dto.Ticket {
	return &dto.Ticket{
		TicketID: ticket.TicketID,
		EventID:  ticket.EventID,
	}
}
