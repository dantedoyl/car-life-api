package models

import "time"

type EventPost struct {
	ID          uint64    `json:"id" binding:"required"`
	Text        string    `json:"text" binding:"required"`
	UserID      uint64    `json:"user_id" binding:"required"`
	EventID     uint64    `json:"event_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" binding:"required"`
	Attachments []string  `json:"attachments" binding:"required"`
}

type CreatePostRequest struct {
	Text        string    `json:"text" binding:"required"`
}
