package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *TicketMongoDBRepository) FindTicket(eventID uuid.UUID, ticketID uuid.UUID) *models.Ticket {
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

func (r *TicketMongoDBRepository) UpdateTicketScanData(ticket *models.Ticket) *models.Ticket {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	updateOptions := new(options.FindOneAndUpdateOptions)
	updateOptions.SetReturnDocument(options.After)

	tickets := r.getCollection()
	result := tickets.FindOneAndUpdate(
		findContext,
		bson.M{"ticket_id": ticket.TicketID, "event_id": ticket.EventID},
		bson.M{
			"$set": bson.M{
				"scanned_by": ticket.ScannedBy,
				"scanned_at": ticket.ScannedAt,
			},
		},
		updateOptions,
	)

	updatedTicket := new(models.Ticket)
	err := result.Decode(updatedTicket)
	if err != nil {
		return nil
	}
	return updatedTicket
}
