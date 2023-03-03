package models

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"math/big"
)

type Ticket struct {
	/*************** ticket key ****************/
	EventID  uuid.UUID `bson:"event_id"`
	TicketID uuid.UUID `bson:"ticket_id"`
	/*******************************************/

	/*********** ticket public info ************/
	TokenID      *big.Int `bson:"token_id"`
	ContractAddr string   `bson:"contract_address"`
	/*******************************************/

	/****************** owner ******************/
	AttendantID         uuid.UUID `bson:"attendant_id"`
	AttendantWalletAddr string    `bson:"attendant_addr"`
	/*******************************************/
}

func (ticket *Ticket) GetTicketFingerprint() string {
	ticketFingerprintData := ticket.ContractAddr + "/" + ticket.TokenID.Text(16)
	h := sha256.New()
	h.Write([]byte(ticketFingerprintData))
	return hex.EncodeToString(h.Sum(nil))
}
