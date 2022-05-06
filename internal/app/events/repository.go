package events

import "github.com/dantedoyl/car-life-api/internal/app/models"

type IEventsRepository interface {
	InsertEvent(event *models.Event) error
	GetEventByID(id int64, userID uint64) (*models.Event, error)
	GetEvents(idGt *uint64, idLte *uint64, limit *uint64, query *string, downLeftLongitude *float32, downLeftLatitude *float32, upperRightLongitude *float32, upperRightLatitude *float32) ([]*models.Event, error)
	UpdateEvent(event *models.Event) (*models.Event, error)
	GetEventsUserByStatus(event_id int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.UserCard, error)
	SetUserStatusByEventID(eventID int64, userID int64, status string) error
	GetEventChatID(eventID int64, userID int64) (int64, error)
	SetEventChatID(eventID int64, chatID int64) error
	DeleteUserFromEvent(eventID int64, userID int64) error
}
