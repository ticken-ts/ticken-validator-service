package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"math/big"
	"ticken-validator-service/log"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
)

type TicketSyncer struct {
	/************ repositories *****+******/
	eventRepo     repos.IEventRepository
	ticketRepo    repos.ITicketRepository
	attendantRepo repos.IAttendantRepository
	/**************************************/

	/********* blockchain callers *********/
	pvtbcCaller *pvtbc.Caller
	pubbcCaller pubbc.Caller
	/**************************************/
}

func NewTicketSyncer(pvtbcCaller *pvtbc.Caller, pubbcCaller pubbc.Caller, repoProvider repos.IProvider) ITicketSyncer {
	return &TicketSyncer{
		pvtbcCaller: pvtbcCaller,
		pubbcCaller: pubbcCaller,

		eventRepo:     repoProvider.GetEventRepository(),
		ticketRepo:    repoProvider.GetTicketRepository(),
		attendantRepo: repoProvider.GetAttendantRepository(),
	}
}

func (syncer *TicketSyncer) Sync(eventID uuid.UUID, callerID uuid.UUID) error {
	event := syncer.eventRepo.FindEvent(eventID)
	if event == nil {
		return fmt.Errorf("event with id %s not found", eventID)
	}

	pvtbcEvent, err := syncer.getPvtbcEvent(event.PvtBCChannel, eventID)
	if err != nil {
		return err
	}

	if event.SyncStatus == models.EventSyncing {
		return fmt.Errorf("event is already syncing")
	}
	if event.SyncStatus == models.EventSynced {
		return fmt.Errorf("event is already synced")
	}

	event.SyncStatus = models.EventSyncing
	if err = syncer.eventRepo.UpdateSyncStatus(event); err != nil {
		return fmt.Errorf("failed to start event syncing")
	}

	go syncer.sync(event, pvtbcEvent)
	return nil
}

func (syncer *TicketSyncer) sync(event *models.Event, pvtbcEvent *chainmodels.Event) {
	for _, section := range pvtbcEvent.Sections {
		pvtbcTickets, err := syncer.pvtbcCaller.GetSectionTickets(event.EventID, section.Name)
		if err != nil {
			syncer.revertSyncingStatus(event, err)
			return
		}

		if len(pvtbcTickets) == 0 {
			continue
		}

		var tickets []*models.Ticket
		for _, pvtbcTicket := range pvtbcTickets {
			syncedTicket, err := syncer.syncTicket(event.PubBCAddress, pvtbcTicket)
			if err != nil {
				syncer.revertSyncingStatus(event, err)
				return
			}
			tickets = append(tickets, syncedTicket)
		}

		if err := syncer.ticketRepo.AddManyTickets(tickets); err != nil {
			syncer.revertSyncingStatus(event, err)
			return
		}
	}

	event.SyncStatus = models.EventSynced
	if err := syncer.eventRepo.UpdateSyncStatus(event); err != nil {
		log.TickenLogger.Error().Msg("failed to update status to synced")
	}
}

func (syncer *TicketSyncer) syncTicket(pubbcAddr string, pvtbcTicket *chainmodels.Ticket) (*models.Ticket, error) {
	tokenID, ok := big.NewInt(0).SetString(pvtbcTicket.TokenID, 16)
	if !ok {
		return nil, fmt.Errorf("failed to obtain token ID")
	}

	pubbcTicket, err := syncer.pubbcCaller.GetTicket(pubbcAddr, tokenID)
	if err != nil {
		return nil, err
	}

	attendant := syncer.attendantRepo.FindAttendantByWalletAddr(pubbcTicket.Owner)
	if attendant == nil {
		return nil, fmt.Errorf("attendant not found")
	}

	ticket := &models.Ticket{
		TicketID: pvtbcTicket.TicketID,
		EventID:  pvtbcTicket.EventID,

		TokenID:      tokenID,
		ContractAddr: pubbcAddr,

		AttendantWalletAddr: pubbcTicket.Owner,
		AttendantID:         attendant.AttendantID,
	}

	return ticket, nil
}

func (syncer *TicketSyncer) getPvtbcEvent(channel string, eventID uuid.UUID) (*chainmodels.Event, error) {
	if err := syncer.pvtbcCaller.SetChannel(channel); err != nil {
		return nil, err
	}
	pvtbcEvent, err := syncer.pvtbcCaller.GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	return pvtbcEvent, nil
}

func (syncer *TicketSyncer) revertSyncingStatus(event *models.Event, err error) {
	event.SyncStatus = models.EventDesynced
	_ = syncer.eventRepo.UpdateSyncStatus(event)
	log.TickenLogger.Error().Msg(err.Error())
}
