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

func (r *AttendantMongoDBRepository) AddAttendant(attendant *models.Attendant) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	attendantsCollection := r.getCollection()
	_, err := attendantsCollection.InsertOne(storeContext, attendant)
	if err != nil {
		return err
	}

	return nil
}

func (r *AttendantMongoDBRepository) FindAttendant(attendantID uuid.UUID) *models.Attendant {
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

func (r *AttendantMongoDBRepository) FindAttendantByWalletAddr(wallerAddr string) *models.Attendant {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	attendantsCollection := r.getCollection()
	result := attendantsCollection.FindOne(findContext, bson.M{"wallet_address": wallerAddr})

	var foundAttendant models.Attendant
	err := result.Decode(&foundAttendant)

	if err != nil {
		return nil
	}

	return &foundAttendant
}

func (r *AttendantMongoDBRepository) AnyWithID(attendantID uuid.UUID) bool {
	return r.FindAttendant(attendantID) != nil
}
