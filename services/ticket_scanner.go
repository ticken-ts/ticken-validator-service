package services

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/google/uuid"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
	"ticken-validator-service/utils"
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

func (s *TicketScanner) Scan(eventID, ticketID uuid.UUID, signature []byte, validatorID uuid.UUID) (*models.Ticket, error) {
	ticket := s.ticketRepo.FindTicket(eventID, ticketID)
	if ticket == nil {
		return nil, fmt.Errorf("ticket not found")
	}

	ticketOwner := s.attendantRepo.FindAttendant(ticket.AttendantID)
	if ticketOwner == nil {
		return nil, fmt.Errorf("attendant not found")
	}

	if err := ticket.Scan(validatorID); err != nil {
		return nil, err
	}

	if err := s.validateSignature(ticketOwner.PublicKey, signature, ticket); err != nil {
		return nil, fmt.Errorf("ticket signature is not valid")
	}

	if ticket := s.ticketRepo.UpdateTicketScanData(ticket); ticket == nil {
		return nil, fmt.Errorf("failed to update ticket scan data")
	}

	return ticket, nil
}

func (s *TicketScanner) validateSignature(publicKey []byte, signature []byte, ticket *models.Ticket) error {
	ticketFingerprint := ticket.GetTicketFingerprint()

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return fmt.Errorf("failed to decode public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("failed to decode public key")
	}

	err = rsa.VerifyPSS(pubKey, crypto.SHA256, utils.HashSHA256(ticketFingerprint), signature, nil)
	if err != nil {
		return err
	}

	return nil
}
