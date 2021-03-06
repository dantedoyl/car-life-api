package models

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	Value     string    `json:"value" binding:"required"`
	UserID    uint64    `json:"user_id" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}

type SignUpResponse struct {
	CarID   int64    `json:"car_id" binding:"required"`
	Session *Session `json:"session" binding:"required"`
}

type User struct {
	VKID        uint64   `json:"vkid" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Surname     string   `json:"surname" binding:"required"`
	AvatarUrl   string   `json:"avatar_url" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
	Description string   `json:"description" binding:"required"`
	CarID       int64    `json:"-"`
}

type UserCard struct {
	VKID      uint64 `json:"vkid" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	AvatarUrl string `json:"avatar_url" binding:"required"`
}

type CarCard struct {
	ID          uint64    `json:"id" binding:"required"`
	AvatarUrl   string    `json:"avatar_url" binding:"required"`
	Brand       string    `json:"brand" binding:"required"`
	Model       string    `json:"model" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Owner     	UserCard  `json:"owner" binding:"required"`
	Body        string    `json:"body" binding:"required"`
	Engine      string    `json:"engine" binding:"required"`
	HorsePower  string    `json:"horse_power" binding:"required"`
	Name        string    `json:"name" binding:"required"`
}

type CarRequest struct {
	Brand       string
	Model       string
	Date        time.Time
	Description string
	Body        string
	Engine      string
	HorsePower  string
	Name        string
}

type SignUpRequest struct {
	VKID        uint64
	Name        string
	Surname     string
	AvatarUrl   string
	Garage      []CarRequest
	Tags        []string
	Description string
}

type UpdateRequest struct {
	Tags        []string
	Description string
}

type LoginRequest struct {
	VKID uint64
}

func CreateSession(userID uint64) *Session {
	return &Session{
		Value:     uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}
}

type Complaint struct {
	UserID int64
	Text string
	TargetID int64
}

type ComplaintReq struct {
	Text string
}