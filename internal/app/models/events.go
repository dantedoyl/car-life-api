package models

import "time"

type Event struct {
	ID             uint64    `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Club   Club `json:"club" binding:"required"`
	Description string `json:"description" binding:"required"`
	EventDate time.Time `json:"event_date" binding:"required"`
	Latitude  float32 `json:"latitude" binding:"required"`
	Longitude float32 `json:"longitude" binding:"required"`
	AvatarUrl    string `json:"avatar" binding:"required"`
}
