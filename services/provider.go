package services

import (
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-validator-service/repos"
)

type Provider struct {
	TicketScanner TicketScanner
	EventManager  EventManager
}

func NewProvider(repoProvider repos.IProvider, pvtbcCaller *pvtbc.Caller) (IProvider, error) {
	provider := new(Provider)

	eventRepo := repoProvider.GetEventRepository()
	ticketRepo := repoProvider.GetTicketRepository()

	provider.EventManager = NewEventManager(eventRepo)
	provider.TicketScanner = NewTicketScanner(eventRepo, ticketRepo, pvtbcCaller)

	return provider, nil
}

func (provider *Provider) GetTicketScanner() TicketScanner {
	return provider.TicketScanner
}

func (provider *Provider) GetEventManager() EventManager {
	return provider.EventManager
}
