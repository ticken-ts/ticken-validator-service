package models

import "strings"

type Ticket struct {
	EventID  string `json:"event_id"`
	TicketID string `json:"ticket_id"`
}

func NewTicket(ticketID string, eventID string) *Ticket {
	return &Ticket{
		EventID:  strings.TrimSpace(eventID),
		TicketID: strings.TrimSpace(ticketID),
	}
}
