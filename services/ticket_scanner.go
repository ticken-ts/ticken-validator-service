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

func (s *TicketScanner) Scan(eventID, ticketID uuid.UUID, signature string, validatorID uuid.UUID) (*models.Ticket, error) {
	ticket := s.ticketRepo.FindTicket(eventID, ticketID)
	if ticket == nil {
		return nil, fmt.Errorf("ticket not found")
	}

	ticketOwner := s.attendantRepo.FindAttendant(ticket.AttendantID)
	if ticketOwner == nil {
		return nil, fmt.Errorf("attendant not found")
	}

	if err := s.validateSignature(ticketOwner.PublicKey, signature, ticket); err != nil {
		return nil, fmt.Errorf("ticket signature is not valid")
	}

	return ticket, nil
}

func (s *TicketScanner) validateSignature(publicKey []byte, signature string, ticket *models.Ticket) error {
	ticketFingerprint := ticket.GetTicketFingerprint()

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return fmt.Errorf("failed to decode public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	if pubKey, ok := pub.(*rsa.PublicKey); ok {
		err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, []byte(ticketFingerprint), []byte(signature))
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("failed to generate public key")
}
