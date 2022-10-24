package models

type Event struct {
	EventID      string `json:"event_id"`
	OrganizerID  string `json:"organizer_id"`
	PvtBCChannel string `json:"pvt_bc_channel"`
}

func NewEvent(EventID string, OrganizerID string, PvtBCChannel string) *Event {
	return &Event{
		EventID:      EventID,
		OrganizerID:  OrganizerID,
		PvtBCChannel: PvtBCChannel,
	}
}
