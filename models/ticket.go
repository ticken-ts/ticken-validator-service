package models

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"ticken-validator-service/utils"
	"time"
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

	/******************* scan ******************/
	ScannedAt time.Time `bson:"scanned_at"`
	ScannedBy uuid.UUID `bson:"scanned_by"`
	/*******************************************/

}

func (ticket *Ticket) GetTicketFingerprint() string {
	ticketFingerprintData := ticket.ContractAddr + "/" + ticket.TokenID.Text(16)
	hash := utils.HashSHA256(ticketFingerprintData)
	return base64.URLEncoding.EncodeToString(hash)
}

func (ticket *Ticket) Scan(validatorID uuid.UUID) error {
	if ticket.ScannedBy != uuid.Nil {
		return fmt.Errorf("ticket already scanned")
	}
	ticket.ScannedBy = validatorID
	ticket.ScannedAt = time.Now()
	return nil
}
