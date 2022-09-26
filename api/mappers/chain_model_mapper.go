package mappers

import (
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-validator-service/api/dto"
)

func MapChainTicketToDTO(ticket *chain_models.Ticket) *dto.Ticket {
	return &dto.Ticket{
		TicketID: ticket.TicketID,
		EventID:  ticket.EventID,
		Status:   ticket.Status,
		Owner:    ticket.Owner,
	}
}
