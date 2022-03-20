package events

import (
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type IClubsUsecase interface {
	CreateClub(event *models.Club) error
	GetClubByID(id uint64) (*models.Club, error)
	GetClubs() ([]*models.Club, error)
	UpdateAvatar(eventID int64, fileHeader *multipart.FileHeader) (*models.Club, error)
}
