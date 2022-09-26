package repos

import "ticken-validator-service/models"

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID string) *models.Event
}

type Provider interface {
	GetEventRepository() EventRepository
}

type Factory interface {
	BuildEventRepository() any
}
