package models

import "github.com/google/uuid"

type Attendant struct {
	AttendantID   uuid.UUID `bson:"attendant_id"`
	WalletAddress string    `bson:"wallet_address"`
	PublicKey     []byte    `bson:"public_key"`
}
