package events

import (
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type IClubsUsecase interface {
	CreateClub(event *models.Club) error
	GetClubByID(id uint64) (*models.Club, error)
	GetClubs(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Club, error)
	UpdateAvatar(eventID int64, fileHeader *multipart.FileHeader) (*models.Club, error)
	GetTags() ([]models.Tag, error)
}
