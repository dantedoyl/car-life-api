package usecase

import (
	"github.com/dantedoyl/car-life-api/internal/app/events"
	"github.com/dantedoyl/car-life-api/internal/app/models"
)

type EventsUsecase struct {
	eventsRepo events.IEventsRepository
}

func NewEventsUsecase(repo events.IEventsRepository) events.IEventsUsecase {
	return &EventsUsecase{
		eventsRepo: repo,
	}
}

func (eu *EventsUsecase) CreateEvent(event *models.Event) error {
	return eu.eventsRepo.InsertEvent(event)
}

func (eu *EventsUsecase) GetEventByID(id uint64) (*models.Event, error) {
	return eu.eventsRepo.GetEventByID(int64(id))
}

func (eu *EventsUsecase) GetEvents() ([]*models.Event, error) {
	return eu.eventsRepo.GetEvents()
}