package models

import "github.com/google/uuid"

type Validator struct {
	ValidatorID uuid.UUID `bson:"validator_id"`
	Firstname   string    `bson:"firstname"`
	Lastname    string    `bson:"lastname"`
	Username    string    `bson:"username"`
	Email       string    `bson:"email"`
}
