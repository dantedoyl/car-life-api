package models

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	Value     string
	UserID    uint64
	ExpiresAt time.Time
}

type User struct {
	VKID uint64
	Name string
	Surname string
	AvatarUrl string
	Garage []*CarCard
	OwnClubs []ClubCard
	Tags []string
	ParticipantClubs []ClubCard
	ParticipantEvents []EventCard
}

type UserCard struct {
	ID uint64
	Name string
}

type CarCard struct {
	ID uint64
	AvatarUrl string
	Barnd string
	Model string
	Date time.Time
	Description string
	OwnerID uint64
}

type CarRequest struct {
	Barnd string
	Model string
	Date time.Time
	Description string
}

type SignUpRequest struct {
	VKID uint64
	Name string
	Surname string
	AvatarUrl string
	Garage []CarRequest
	Tags []string
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
