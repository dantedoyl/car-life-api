package events

import "github.com/dantedoyl/car-life-api/internal/app/models"

type IEventsRepository interface {
	InsertEvent(event *models.Event) error
	GetEventByID(id int64) (*models.Event, error)
	GetEvents(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Event, error)
	UpdateEvent(event *models.Event) (*models.Event, error)
}
