package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-validator-service/config"
	"ticken-validator-service/infra"
)

type Factory struct {
	dbClient *mongo.Client
	dbName   string
}

func NewFactory(db infra.Db, dbConfig *config.DatabaseConfig) *Factory {
	return &Factory{
		dbClient: db.GetClient().(*mongo.Client),
		dbName:   dbConfig.Name,
	}
}

func (factory *Factory) BuildEventRepository() any {
	return NewEventRepository(factory.dbClient, factory.dbName)
}

func (factory *Factory) BuildTicketRepository() any {
	return NewTicketRepository(factory.dbClient, factory.dbName)
}
