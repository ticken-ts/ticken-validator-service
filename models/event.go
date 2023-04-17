package models

import "github.com/google/uuid"

type EventSyncStatus string

const (
	EventDesynced EventSyncStatus = "desynced"
	EventSynced   EventSyncStatus = "synced"
	EventSyncing  EventSyncStatus = "syncing"
)

type Event struct {
	EventID        uuid.UUID       `bson:"event_id"`
	OrganizerID    uuid.UUID       `bson:"organizer_id"`
	PvtBCChannel   string          `bson:"pvt_bc_channel"`
	PubBCAddress   string          `bson:"pub_bc_address"`
	OrganizationID uuid.UUID       `bson:"organization_id"`
	SyncStatus     EventSyncStatus `bson:"sync_status"`
}
