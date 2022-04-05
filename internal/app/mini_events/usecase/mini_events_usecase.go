package usecase

import (
	"github.com/dantedoyl/car-life-api/internal/app/mini_events"
	"github.com/dantedoyl/car-life-api/internal/app/models"
)

type MiniEventsUsecase struct {
	miniEventsRepo mini_events.IMiniEventsRepository
}

func NewMiniEventsUsecase(repo mini_events.IMiniEventsRepository) mini_events.IMiniEventsUsecase {
	return &MiniEventsUsecase{
		miniEventsRepo: repo,
	}
}

func (mu *MiniEventsUsecase) CreateMiniEvent(event *models.MiniEvent) error {
	return mu.miniEventsRepo.InsertMiniEvent(event)
}

func (mu *MiniEventsUsecase) GetMiniEventByID(id uint64) (*models.MiniEvent, error) {
	return mu.miniEventsRepo.GetMiniEventByID(int64(id))
}

func (mu *MiniEventsUsecase) GetMiniEvents(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.MiniEvent, error) {
	return mu.miniEventsRepo.GetMiniEvents(idGt, idLte, limit, query)
}
