package usecase

import (
	"github.com/dantedoyl/car-life-api/internal/app/clients/filesystem"
	"github.com/dantedoyl/car-life-api/internal/app/events"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
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

func (eu *EventsUsecase) GetEventByID(id uint64, userID uint64) (*models.Event, error) {
	return eu.eventsRepo.GetEventByID(int64(id), userID)
}

func (eu *EventsUsecase) GetEvents(idGt *uint64, idLte *uint64, limit *uint64, query *string, downLeftLongitude *float32, downLeftLatitude *float32, upperRightLongitude *float32, upperRightLatitude *float32) ([]*models.Event, error) {
	return eu.eventsRepo.GetEvents(idGt, idLte, limit, query, downLeftLongitude, downLeftLatitude, upperRightLongitude, upperRightLatitude)
}

func (eu *EventsUsecase) UpdateAvatar(eventID int64, fileHeader *multipart.FileHeader) (*models.Event, error) {
	event, err := eu.eventsRepo.GetEventByID(eventID, 0)
	if err != nil {
		return nil, err
	}

	imgUrl, err := filesystem.InsertPhoto(fileHeader, "img/events/")
	if err != nil {
		return nil, err
	}

	oldAvatar := event.AvatarUrl
	event.AvatarUrl = imgUrl
	event, err = eu.eventsRepo.UpdateEvent(event)
	if err != nil {
		return nil, err
	}

	if oldAvatar == "/img/events/default.webp" {
		return event, nil
	}

	err = filesystem.RemovePhoto(oldAvatar)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (eu *EventsUsecase) GetEventsUserByStatus(event_id int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.UserCard, error) {
	return eu.eventsRepo.GetEventsUserByStatus(event_id, status, idGt, idLte, limit)
}

func (eu *EventsUsecase) SetUserStatusByEventID(eventID int64, userID int64, status string) error {
	return eu.eventsRepo.SetUserStatusByEventID(eventID, userID, status)
}

func (eu *EventsUsecase) ApproveRejectUserParticipateInEvent(eventID int64, userID int64, decision string) error {
	if decision == "approve" {
		return eu.eventsRepo.SetUserStatusByEventID(eventID, userID, "participant")
	}
	return eu.eventsRepo.SetUserStatusByEventID(eventID, userID, "spectator")
}

func (eu *EventsUsecase) GetEventChatID(eventID int64, userID int64) (int64, error) {
	return eu.eventsRepo.GetEventChatID(eventID, userID)
}

func (eu *EventsUsecase) SetEventChatID(eventID int64, chatID int64) error {
	return eu.eventsRepo.SetEventChatID(eventID, chatID)
}

func (eu *EventsUsecase) DeleteUserFromEvent(eventID int64, chatID int64) error {
	return eu.eventsRepo.DeleteUserFromEvent(eventID, chatID)
}

func (eu *EventsUsecase) DeleteEventByID(eventID int64) error {
	return eu.eventsRepo.DeleteEventByID(eventID)
}

func (eu *EventsUsecase) ComplainByID(complaint models.Complaint) error {
	return eu.eventsRepo.ComplainByID(complaint)
}
