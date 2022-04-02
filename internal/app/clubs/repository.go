package events

import "github.com/dantedoyl/car-life-api/internal/app/models"

type IClubsRepository interface {
	InsertClub(event *models.Club) error
	GetClubByID(id int64, userID uint64) (*models.Club, error)
	GetClubs(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Club, error)
	UpdateClub(event *models.Club) (*models.Club, error)
	GetTags() ([]models.Tag, error)
	GetClubsUserByStatus(club_id int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.UserCard, error)
	GetClubsCars(club_id int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.CarCard, error)
	GetClubsEvents(club_id int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventCard, error)
	SetUserStatusByClubID(clubID int64, userID int64, status string) error
	GetUserStatusInClub(clubID int64, userID int64) (*models.ClubUser, error)
}
