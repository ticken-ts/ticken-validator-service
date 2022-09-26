package services

import (
	"fmt"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"ticken-validator-service/repos"
)

type ticketScanner struct {
	pvtbcCaller     *pvtbc.Caller
	eventRepository repos.EventRepository
}

func NewTicketScanner(eventRepository repos.EventRepository, pvtbcCaller *pvtbc.Caller) TicketScanner {
	return &ticketScanner{
		pvtbcCaller:     pvtbcCaller,
		eventRepository: eventRepository,
	}
}

func (s *ticketScanner) Scan(eventID string, ticketID string, owner string) (*chain_models.Ticket, error) {
	event := s.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	err := s.pvtbcCaller.SetChannel(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	ticket, err := s.pvtbcCaller.ScanTicket(
		ticketID,
		eventID,
		owner,
	)

	if err != nil {
		return nil, err
	}

	return ticket, nil
}
