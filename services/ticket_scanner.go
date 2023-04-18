package services

import (
	"crypto/ecdsa"
	"encoding/pem"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"math/big"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
)

type TicketScanner struct {
	eventRepo     repos.IEventRepository
	ticketRepo    repos.ITicketRepository
	attendantRepo repos.IAttendantRepository
}

func NewTicketScanner(repoProvider repos.IProvider) *TicketScanner {
	return &TicketScanner{
		eventRepo:     repoProvider.GetEventRepository(),
		ticketRepo:    repoProvider.GetTicketRepository(),
		attendantRepo: repoProvider.GetAttendantRepository(),
	}
}

func (scanner *TicketScanner) Scan(eventID, ticketID, validatorID uuid.UUID, rSignatureField, sSignatureField string) (*models.Ticket, error) {
	ticket := scanner.ticketRepo.FindTicket(eventID, ticketID)
	if ticket == nil {
		return nil, fmt.Errorf("ticket not found")
	}

	ticketOwner := scanner.attendantRepo.FindAttendant(ticket.AttendantID)
	if ticketOwner == nil {
		return nil, fmt.Errorf("attendant not found")
	}

	if err := ticket.Scan(validatorID); err != nil {
		return nil, err
	}

	publicKey, err := parsePublicKey(ticketOwner.PemPublicKey)
	if err != nil {
		return nil, err
	}

	if err := scanner.validateSignature(publicKey, rSignatureField, sSignatureField, ticket); err != nil {
		return nil, fmt.Errorf("ticket signature is not valid")
	}

	if ticket := scanner.ticketRepo.UpdateTicketScanData(ticket); ticket == nil {
		return nil, fmt.Errorf("failed to update ticket scan data")
	}

	return ticket, nil
}

func (scanner *TicketScanner) validateSignature(publicKey *ecdsa.PublicKey, rSignatureField, sSignatureField string, ticket *models.Ticket) error {
	ticketFingerprint := ticket.GetTicketFingerprint()

	r, ok := big.NewInt(0).SetString(rSignatureField, 16)
	if !ok {
		return fmt.Errorf("failed to read R signature filed")
	}
	s, ok := big.NewInt(0).SetString(sSignatureField, 16)
	if !ok {
		return fmt.Errorf("failed to read S signature filed")
	}

	signature := make([]byte, 0)
	signature = append(signature, r.Bytes()...)
	signature = append(signature, s.Bytes()...)

	pubKey := crypto.CompressPubkey(publicKey)

	if ok := crypto.VerifySignature(pubKey, ticketFingerprint, signature); !ok {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

func parsePublicKey(pemEncodedPublicKey string) (*ecdsa.PublicKey, error) {
	pemPublicKey, _ := pem.Decode([]byte(pemEncodedPublicKey))
	pubKey, err := crypto.UnmarshalPubkey(pemPublicKey.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode pem public key: %s", err.Error())
	}
	return pubKey, nil
}
