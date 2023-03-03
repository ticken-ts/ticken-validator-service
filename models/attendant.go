package models

import "github.com/google/uuid"

type Attendant struct {
	AttendantID uuid.UUID `bson:"attendant_id"`
	PublicKey   []byte    `bson:"public_key"`
}
