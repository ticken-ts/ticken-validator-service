package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"math/big"
	"ticken-validator-service/env"
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

func (ticket *Ticket) GetTicketFingerprint() []byte {
	ticketFingerprintData := ticket.ContractAddr + "/" + ticket.TokenID.Text(16) + "/" + ticket.getTOTPToken()
	return utils.HashSHA256(ticketFingerprintData)
}

func (ticket *Ticket) Scan(validatorID uuid.UUID) error {
	if ticket.ScannedBy != uuid.Nil {
		return fmt.Errorf("ticket already scanned")
	}
	ticket.ScannedBy = validatorID
	ticket.ScannedAt = time.Now()
	return nil
}

func (ticket *Ticket) getTOTPToken() string {
	token, _ := totp.GenerateCodeCustom(env.TickenEnv.TOPTSecret, time.Now(), totp.ValidateOpts{
		Period:    60 * 30,
		Algorithm: otp.AlgorithmSHA512,
		Digits:    8,
	})
	return token
}
