package services

import (
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
)

type ticketScanner struct {
	pvtbcCaller      *pvtbc.Caller
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
}

func NewTicketScanner(eventRepository repos.EventRepository, ticketRepository repos.TicketRepository, pvtbcCaller *pvtbc.Caller) TicketScanner {
	return &ticketScanner{
		pvtbcCaller:      pvtbcCaller,
		eventRepository:  eventRepository,
		ticketRepository: ticketRepository,
	}
}

func (s *ticketScanner) Scan(eventID string, ticketID string, owner string) (*models.Ticket, error) {
	//event := s.eventRepository.FindEvent(eventID)
	//if event == nil {
	//	return nil, fmt.Errorf("could not determine organizer channel")
	//}
	//
	//if err := s.pvtbcCaller.SetChannel(event.PvtBCChannel); err != nil {
	//	return nil, err
	//}
	//
	//pvtbcTicket, err := s.pvtbcCaller.ScanTicket(ticketID, eventID, owner)
	//if err != nil {
	//	return nil, err
	//}
	//
	//ticket := models.NewTicket(pvtbcTicket.TicketID, pvtbcTicket.EventID)
	//if err := s.ticketRepository.AddTicket(ticket); err != nil {
	//	return nil, err
	//}
	//
	//return ticket, nil
	return nil, nil
}
