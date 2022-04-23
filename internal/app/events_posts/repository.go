package events_posts

import "github.com/dantedoyl/car-life-api/internal/app/models"

type IEventsPostsRepository interface {
	InsertEventPost(event *models.EventPost) error
	GetEventPostByPostID (postID uint64) (*models.EventPost, error)
	GetEventsPostsByEventID(eventID uint64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventPost, error)
	InsertEventPostAttachments(postID uint64, attachments []string) error
}
