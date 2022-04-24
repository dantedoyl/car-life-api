package models

import "time"

type EventPost struct {
	ID          uint64    `json:"id" binding:"required"`
	Text        string    `json:"text" binding:"required"`
	User        UserCard  `json:"user" binding:"required"`
	EventID     uint64    `json:"event_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" binding:"required"`
	Attachments []string  `json:"attachments" binding:"required"`
}

type CreatePostRequest struct {
	Text        string    `json:"text" binding:"required"`
}
