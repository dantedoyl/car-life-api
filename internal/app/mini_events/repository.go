package mini_events

import "github.com/dantedoyl/car-life-api/internal/app/models"

type IMiniEventsRepository interface {
	InsertMiniEvent(event *models.MiniEvent) error
	GetMiniEventByID(id int64) (*models.MiniEvent, error)
	GetMiniEvents(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.MiniEvent, error)
}
