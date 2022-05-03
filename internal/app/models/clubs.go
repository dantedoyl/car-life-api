package models

type Club struct {
	ID                uint64      `json:"id" binding:"required"`
	Name              string      `json:"name" binding:"required"`
	Description       string      `json:"description" binding:"required"`
	AvatarUrl         string      `json:"avatar" binding:"required"`
	Tags              []string    `json:"tags" binding:"required"`
	EventsCount       int         `json:"events_count" binding:"required"`
	ParticipantsCount int         `json:"participants_count" binding:"required"`
	SubscribersCount int          `json:"subscribers_count" binding:"required"`
	OwnerID           uint64      `json:"owner_id" binding:"required"`
	UserStatus string `json:"user_status" binding:"required"`
}

type ClubUser struct {
	UserID int64
	ClubID int64
	Status string
}

type ClubCard struct {
	ID                uint64   `json:"id" binding:"required"`
	Name              string   `json:"name" binding:"required"`
	AvatarUrl         string   `json:"avatar" binding:"required"`
	Tags              []string `json:"tags" binding:"required"`
	ParticipantsCount int      `json:"participants_count" binding:"required"`
	SubscribersCount int          `json:"subscribers_count" binding:"required"`
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

type ChatLink struct {
	ChatLink string `json:"chat_link" binding:"required"`
}
