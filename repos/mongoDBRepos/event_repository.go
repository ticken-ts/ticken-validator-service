package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-validator-service/models"
)

const EventCollectionName = "events"

type EventMongoDBRepository struct {
	baseRepository
}

func NewEventRepository(dbClient *mongo.Client, dbName string) *EventMongoDBRepository {
	return &EventMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: EventCollectionName,
		},
	}
}

func (r *EventMongoDBRepository) AddEvent(event *models.Event) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	_, err := events.InsertOne(storeContext, event)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventMongoDBRepository) FindEvent(eventID uuid.UUID) *models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result := events.FindOne(findContext, bson.M{"event_id": eventID})

	var foundEvent models.Event
	err := result.Decode(&foundEvent)

	if err != nil {
		return nil
	}

	return &foundEvent
}

func (r *EventMongoDBRepository) AnyWithID(eventID uuid.UUID) bool {
	return r.FindEvent(eventID) != nil
}
