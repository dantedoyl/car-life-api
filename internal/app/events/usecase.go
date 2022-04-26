package events

import (
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type IEventsUsecase interface {
	CreateEvent(event *models.Event) error
	GetEventByID(id uint64, userID uint64) (*models.Event, error)
	GetEvents(idGt *uint64, idLte *uint64, limit *uint64, query *string, onlyActual bool,downLeftLongitude *float32, downLeftLatitude *float32, upperRightLongitude *float32, upperRightLatitude *float32) ([]*models.Event, error)
	UpdateAvatar(eventID int64, fileHeader *multipart.FileHeader) (*models.Event, error)
	GetEventsUserByStatus(event_id int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.UserCard, error)
	SetUserStatusByEventID(eventID int64, userID int64, status string) error
	ApproveRejectUserParticipateInEvent(eventID int64, userID int64, decision string) error
	GetEventChatID(eventID int64, userID int64) (int64, error)
	SetEventChatID(eventID int64, chatID int64) error
}
