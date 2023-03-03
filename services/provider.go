package services

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-validator-service/repos"
)

type Provider struct {
	TicketScanner ITicketScanner
	TicketSyncer  ITicketSyncer
	EventManager  IEventManager
}

func NewProvider(repoProvider repos.IProvider, pvtbcCaller *pvtbc.Caller, pubbcCaller pubbc.Caller) (IProvider, error) {
	provider := new(Provider)

	provider.TicketSyncer = NewTicketSyncer(pvtbcCaller, pubbcCaller, repoProvider)
	provider.TicketScanner = NewTicketScanner(repoProvider)
	provider.EventManager = NewEventManager(repoProvider)

	return provider, nil
}

func (provider *Provider) GetTicketScanner() ITicketScanner {
	return provider.TicketScanner
}

func (provider *Provider) GetTicketSyncer() ITicketSyncer {
	return provider.TicketSyncer
}

func (provider *Provider) GetEventManager() IEventManager {
	return provider.EventManager
}
