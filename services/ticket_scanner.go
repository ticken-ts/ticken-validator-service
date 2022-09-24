package services

import (
	"fmt"
	"ticken-validator-service/blockchain/pvtbc"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
)

type ticketScanner struct {
	pvtbcConnector  pvtbc.TickenConnector
	eventRepository repos.EventRepository
}

func NewTicketScanner(eventRepository repos.EventRepository, pvtbcConnector pvtbc.TickenConnector) TicketScanner {
	return &ticketScanner{
		pvtbcConnector:  pvtbcConnector,
		eventRepository: eventRepository,
	}
}

func (s *ticketScanner) Scan(eventID string, ticketID string, owner string) (*models.Ticket, error) {
	event := s.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	err := s.pvtbcConnector.Connect(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	ticketResponse, err := s.pvtbcConnector.ScanTicket(
		ticketID,
		eventID,
		owner,
	)

	if err != nil {
		return nil, err
	}

	ticket := new(models.Ticket)
	ticket.EventID = eventID
	ticket.TicketID = ticketID
	ticket.Status = ticketResponse.Status

	return ticket, nil
}
