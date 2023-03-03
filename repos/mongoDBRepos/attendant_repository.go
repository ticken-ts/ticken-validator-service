package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-validator-service/models"
)

const AttendantCollectionName = "attendants"

type AttendantMongoDBRepository struct {
	baseRepository
}

func NewAttendantRepository(dbClient *mongo.Client, dbName string) *AttendantMongoDBRepository {
	return &AttendantMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: AttendantCollectionName,
		},
	}
}

func (r *EventMongoDBRepository) AddAttendant(attendant *models.Attendant) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	attendantsCollection := r.getCollection()
	_, err := attendantsCollection.InsertOne(storeContext, attendant)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventMongoDBRepository) FindAttendant(attendantID uuid.UUID) *models.Attendant {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	attendantsCollection := r.getCollection()
	result := attendantsCollection.FindOne(findContext, bson.M{"attendant_id": attendantID})

	var foundAttendant models.Attendant
	err := result.Decode(&foundAttendant)

	if err != nil {
		return nil
	}

	return &foundAttendant
}

func (r *EventMongoDBRepository) FindAttendantByWalletAddr(wallerAddr string) *models.Attendant {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	attendantsCollection := r.getCollection()
	result := attendantsCollection.FindOne(findContext, bson.M{"wallet_addr": wallerAddr})

	var foundAttendant models.Attendant
	err := result.Decode(&foundAttendant)

	if err != nil {
		return nil
	}

	return &foundAttendant
}

func (r *EventMongoDBRepository) AddManyTickets(tickets []*models.Ticket) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	attendantsCollection := r.getCollection()

	toAdd := make([]interface{}, len(tickets))
	for i, ticket := range tickets {
		toAdd[i] = ticket
	}

	_, err := attendantsCollection.InsertMany(storeContext, toAdd)
	if err != nil {
		return err
	}

	return nil
}
