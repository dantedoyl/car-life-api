package models

type Club struct {
	ID             uint64    `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	AvatarUrl    string `json:"avatar" binding:"required"`
	EventsCount int `json:"events_count" binding:"required"`
	ParticipantsCount int `json:"participants_count" binding:"required"`
}
