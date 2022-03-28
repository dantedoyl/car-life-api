package events

import "github.com/dantedoyl/car-life-api/internal/app/models"

type IClubsRepository interface {
	InsertClub(event *models.Club) error
	GetClubByID(id int64) (*models.Club, error)
	GetClubs(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Club, error)
	UpdateClub(event *models.Club) (*models.Club, error)
	GetTags() ([]models.Tag, error)
}
