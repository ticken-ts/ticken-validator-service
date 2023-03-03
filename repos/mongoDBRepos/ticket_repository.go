package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-validator-service/models"
)

const TicketCollectionName = "tickets"

type TicketMongoDBRepository struct {
	baseRepository
}

func NewTicketRepository(db *mongo.Client, database string) *TicketMongoDBRepository {
	return &TicketMongoDBRepository{
		baseRepository{
			dbClient:       db,
			dbName:         database,
			collectionName: TicketCollectionName,
		},
	}
}

func (r *TicketMongoDBRepository) AddTicket(ticket *models.Ticket) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	tickets := r.getCollection()
	_, err := tickets.InsertOne(storeContext, ticket)
	if err != nil {
		return err
	}

	return nil
}

func (r *TicketMongoDBRepository) AddManyTickets(tickets []*models.Ticket) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	ticketCollection := r.getCollection()

	ticketsData := make([]interface{}, len(tickets))
	for i, ticket := range tickets {
		ticketsData[i] = ticket
	}

	_, err := ticketCollection.InsertMany(storeContext, ticketsData)
	if err != nil {
		return err
	}

	return nil
}

func (r *TicketMongoDBRepository) FindTicket(eventID string, ticketID string) *models.Ticket {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	tickets := r.getCollection()
	result := tickets.FindOne(findContext, bson.M{"event_id": eventID, "ticket_id": ticketID})

	var foundTicket models.Ticket
	err := result.Decode(&foundTicket)

	if err != nil {
		return nil
	}

	return &foundTicket
}
