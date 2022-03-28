package models

type Club struct {
	ID                uint64      `json:"id" binding:"required"`
	Name              string      `json:"name" binding:"required"`
	Description       string      `json:"description" binding:"required"`
	AvatarUrl         string      `json:"avatar" binding:"required"`
	Tags              []string    `json:"tags" binding:"required"`
	EventsCount       int         `json:"events_count" binding:"required"`
	ParticipantsCount int         `json:"participants_count" binding:"required"`
	OwnerID           uint64      `json:"owner_id" binding:"required"`
	ClubGarage        []CarCard   `json:"club_garage" binding:"required"`
	ClubEvents        []EventCard `json:"club_events" binding:"required"`
}

type ClubCard struct {
	ID        uint64   `json:"id" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	AvatarUrl string   `json:"avatar" binding:"required"`
	Tags      []string `json:"tags" binding:"required"`
}

type ClubQuery struct {
	IdGt  *uint64
	IdLte *uint64
	Limit *uint64
	Query *string
}

type CreateClubRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	AvatarUrl   string   `json:"avatar" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
}

type Tag struct {
	ID   uint64 `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
