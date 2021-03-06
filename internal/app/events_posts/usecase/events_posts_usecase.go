package usecase

import (
	"github.com/dantedoyl/car-life-api/internal/app/clients/filesystem"
	"github.com/dantedoyl/car-life-api/internal/app/events_posts"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type EventsPostsUsecase struct {
	eventsPostsRepo events_posts.IEventsPostsRepository
}

func NewEventsPostsUsecase(repo events_posts.IEventsPostsRepository) events_posts.IEventsPostsUsecase {
	return &EventsPostsUsecase{
		eventsPostsRepo: repo,
	}
}

func (epu *EventsPostsUsecase) CreateEventPost(eventPost *models.EventPost) error {
	return epu.eventsPostsRepo.InsertEventPost(eventPost)
}

func (epu *EventsPostsUsecase) GetEventsPostsByEventID(eventID uint64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventPost, error) {
	return epu.eventsPostsRepo.GetEventsPostsByEventID(eventID, idGt, idLte, limit)
}

func (epu *EventsPostsUsecase) UploadAttachments(postID uint64, fileHeader []*multipart.FileHeader) (*models.EventPost, error) {
	event, err := epu.eventsPostsRepo.GetEventPostByPostID(postID)
	if err != nil {
		return nil, err
	}

	imgUrl, err := filesystem.InsertPhotos(fileHeader, "img/events-posts/")
	if err != nil {
		return nil, err
	}

	err = epu.eventsPostsRepo.InsertEventPostAttachments(postID, imgUrl)
	if err != nil {
		return nil, err
	}

	event.Attachments = imgUrl

	return event, nil
}

func (epu *EventsPostsUsecase) GetEventPostByPostID (postID uint64) (*models.EventPost, error) {
	return epu.eventsPostsRepo.GetEventPostByPostID(postID)
}

func (epu *EventsPostsUsecase) DeletePostByID (postID int64) error {
	return epu.eventsPostsRepo.DeletePostByID(postID)
}

func (epu *EventsPostsUsecase) ComplainByID (complaint models.Complaint) error {
	return epu.eventsPostsRepo.ComplainByID(complaint)
}
