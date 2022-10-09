package services

import (
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-validator-service/infra"
	"ticken-validator-service/repos"
	"ticken-validator-service/utils"
)

type provider struct {
	TicketScanner TicketScanner
	EventManager  EventManager
}

func NewProvider(db infra.Db, pvtbcCaller *pvtbc.Caller, tickenConfig *utils.TickenConfig) (Provider, error) {
	provider := new(provider)

	repoProvider, err := repos.NewProvider(db, tickenConfig)
	if err != nil {
		return nil, err
	}

	provider.TicketScanner = NewTicketScanner(repoProvider.GetEventRepository(), repoProvider.GetTicketRepository(), pvtbcCaller)
	provider.EventManager = NewEventManager(repoProvider.GetEventRepository())

	return provider, nil
}

func (provider *provider) GetTicketScanner() TicketScanner {
	return provider.TicketScanner
}

func (provider *provider) GetEventManager() EventManager {
	return provider.EventManager
}
