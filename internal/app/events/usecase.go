package events

import (
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type IEventsUsecase interface {
	CreateEvent(event *models.Event) error
	GetEventByID(id uint64) (*models.Event, error)
	GetEvents() ([]*models.Event, error)
	UpdateAvatar(eventID int64, fileHeader *multipart.FileHeader) (*models.Event, error)
}
