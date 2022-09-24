package services

import (
	"ticken-validator-service/blockchain/pvtbc"
	"ticken-validator-service/infra"
	"ticken-validator-service/repos"
	"ticken-validator-service/utils"
)

type provider struct {
	TicketScanner TicketScanner
}

func NewProvider(db infra.Db, tickenConfig *utils.TickenConfig) (Provider, error) {
	provider := new(provider)

	repoProvider, err := repos.NewProvider(db, tickenConfig)
	if err != nil {
		return nil, err
	}

	pvtbcTickenConnector, err := pvtbc.NewConnector()
	if err != nil {
		return nil, err
	}

	provider.TicketScanner = NewTicketScanner(
		repoProvider.GetEventRepository(),
		pvtbcTickenConnector,
	)

	return provider, nil
}

func (provider *provider) GetTicketScanner() TicketScanner {
	return provider.TicketScanner
}
