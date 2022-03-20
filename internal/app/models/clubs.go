package models

type Club struct {
	ID             uint64    `json:"id"`
	Name     string `json:"name"`
	Description string `json:"description"`
	AvatarUrl    string `json:"avatar"`
	EventsCount int `json:"events_count"`
	ParticipantsCount int `json:"participants_count"`
}
