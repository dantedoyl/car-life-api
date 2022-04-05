package mini_events

import (
	"github.com/dantedoyl/car-life-api/internal/app/models"
)

type IMiniEventsUsecase interface {
	CreateMiniEvent(event *models.MiniEvent) error
	GetMiniEventByID(id uint64) (*models.MiniEvent, error)
	GetMiniEvents(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.MiniEvent, error)
}
