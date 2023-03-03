package services

import (
	"github.com/google/uuid"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
)

type EventManager struct {
	eventRepo repos.IEventRepository
}

func NewEventManager(repoProvider repos.IProvider) *EventManager {
	return &EventManager{eventRepo: repoProvider.GetEventRepository()}
}

func (eventManager *EventManager) AddEvent(eventID, organizerID uuid.UUID, pvtBCChannel, pubBCAddress string) (*models.Event, error) {
	event := &models.Event{
		EventID:      eventID,
		OrganizerID:  organizerID,
		PvtBCChannel: pvtBCChannel,
		PubBCAddress: pubBCAddress,
	}

	if err := eventManager.eventRepo.AddEvent(event); err != nil {
		return nil, err
	}
	return event, nil
}
