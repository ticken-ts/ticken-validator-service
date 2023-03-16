package repos

import (
	"fmt"
	"ticken-validator-service/config"
	"ticken-validator-service/infra"
	"ticken-validator-service/repos/mongoDBRepos"
)

type Provider struct {
	reposFactory        IFactory
	eventRepository     IEventRepository
	ticketRepository    ITicketRepository
	attendantRepository IAttendantRepository
	validatorRepository IValidatorRepository
}

func NewProvider(db infra.Db, dbConfig *config.DatabaseConfig) (*Provider, error) {
	provider := new(Provider)

	switch dbConfig.Driver {
	case config.MongoDriver:
		provider.reposFactory = mongoDBRepos.NewFactory(db, dbConfig)
	default:
		return nil, fmt.Errorf("database driver %s not implemented", dbConfig.Driver)
	}

	return provider, nil
}

func (provider *Provider) GetEventRepository() IEventRepository {
	if provider.eventRepository == nil {
		provider.eventRepository = provider.reposFactory.BuildEventRepository().(IEventRepository)
	}
	return provider.eventRepository
}

func (provider *Provider) GetTicketRepository() ITicketRepository {
	if provider.ticketRepository == nil {
		provider.ticketRepository = provider.reposFactory.BuildTicketRepository().(ITicketRepository)
	}
	return provider.ticketRepository
}

func (provider *Provider) GetAttendantRepository() IAttendantRepository {
	if provider.attendantRepository == nil {
		provider.attendantRepository = provider.reposFactory.BuildAttendantRepository().(IAttendantRepository)
	}
	return provider.attendantRepository
}

func (provider *Provider) GetValidatorRepository() IValidatorRepository {
	if provider.validatorRepository == nil {
		provider.validatorRepository = provider.reposFactory.BuildValidatorRepository().(IValidatorRepository)
	}
	return provider.validatorRepository
}
