package models

import "time"

type Event struct {
	ID             uint64    `json:"id"`
	Name     string `json:"name"`
	Club   Club `json:"club"`
	Description string `json:"description"`
	EventDate time.Time `json:"event_date"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	AvatarUrl    string `json:"avatar"`
}

type Club struct {
	ID             uint64    `json:"id"`
}
