package models

type Club struct {
	ID             uint64    `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	AvatarUrl    string `json:"avatar" binding:"required"`
	Tags        []Tag `json:"tags" binding:"required"`
	EventsCount int `json:"events_count" binding:"required"`
	ParticipantsCount int `json:"participants_count" binding:"required"`
}

type ClubQuery struct {
	IdGt  *uint64
	IdLte *uint64
	Limit *uint64
	Query *string
}

type CreateClubRequest struct {
	Name     string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	AvatarUrl    string `json:"avatar" binding:"required"`
	Tags        []Tag `json:"tags" binding:"required"`
}

type Tag string
