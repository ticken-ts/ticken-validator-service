package mappers

import (
	"ticken-validator-service/api/dto"
	"ticken-validator-service/models"
)

func MapTicketToDTO(ticket *models.Ticket) *dto.TicketDTO {
	return &dto.TicketDTO{
		TicketID: ticket.TicketID.String(),
		EventID:  ticket.EventID.String(),
	}
}
