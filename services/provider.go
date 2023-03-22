package services

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-validator-service/infra"
	"ticken-validator-service/repos"
)

type Provider struct {
	AttendantManager IAttendantManager
	TicketScanner    ITicketScanner
	TicketSyncer     ITicketSyncer
	EventManager     IEventManager
	ValidatorManager IValidatorManager
}

func NewProvider(repoProvider repos.IProvider, builder infra.IBuilder, pvtbcCaller *pvtbc.Caller, pubbcCaller pubbc.Caller) (IProvider, error) {
	provider := new(Provider)

	provider.ValidatorManager = NewValidatorManager(repoProvider, builder.BuildJWTVerifier())
	provider.TicketSyncer = NewTicketSyncer(pvtbcCaller, pubbcCaller, repoProvider)
	provider.AttendantManager = NewAttendantManager(repoProvider)
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

func (provider *Provider) GetAttendantManager() IAttendantManager {
	return provider.AttendantManager
}

func (provider *Provider) GetValidatorManager() IValidatorManager {
	return provider.ValidatorManager
}
