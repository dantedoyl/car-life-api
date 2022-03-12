package events

import "github.com/dantedoyl/car-life-api/internal/app/models"

type IEventsUsecase interface {
	CreateEvent(event *models.Event) error
	GetEventByID(id uint64) (*models.Event, error)
	GetEvents() ([]*models.Event, error)
}
