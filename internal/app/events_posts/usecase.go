package events_posts

import (
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type IEventsPostsUsecase interface {
	CreateEventPost(event *models.EventPost) error
	GetEventsPostsByEventID(eventID uint64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventPost, error)
	UploadAttachments(postID uint64, fileHeader []*multipart.FileHeader) (*models.EventPost, error)
	GetEventPostByPostID (postID uint64) (*models.EventPost, error)
	DeletePostByID(postID int64) error
	ComplainByID(complaint models.Complaint) error
}
