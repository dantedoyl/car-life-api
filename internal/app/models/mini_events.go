package models

import "time"

type MiniEvent struct {
	ID          uint64    `json:"id" binding:"required"`
	Type        MiniEventType `json:"type" binding:"required"`
	User        UserCard  `json:"user" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time `json:"created_at" binding:"required"`
	EndedAt     time.Time `json:"ended_at" binding:"required"`
	Latitude    float32   `json:"latitude" binding:"required"`
	Longitude   float32   `json:"longitude" binding:"required"`
}

type MiniEventType struct {
	ID          uint64    `json:"id" binding:"required"`
	PublicName  string    `json:"public_name" binding:"required"`
	PublicDescription string `json:"public_description" binding:"required"`
}

type CreateMiniEventRequest struct {
	TypeID int64 `json:"type_id" binding:"required"`
	EndedAt time.Time `json:"ended_at" binding:"required"`
	Description string `json:"description" binding:"required"`
	Latitude    float32   `json:"latitude" binding:"required"`
	Longitude   float32   `json:"longitude" binding:"required"`
}
