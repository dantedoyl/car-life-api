package events

import (
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type IClubsUsecase interface {
	CreateClub(event *models.Club) error
	GetClubByID(id uint64, userID uint64) (*models.Club, error)
	GetClubs(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Club, error)
	UpdateAvatar(eventID int64, fileHeader *multipart.FileHeader) (*models.Club, error)
	GetTags() ([]models.Tag, error)
	GetClubsUserByStatus(club_id int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.UserCard, error)
	GetClubsCars(club_id int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.CarCard, error)
	GetClubsEvents(club_id int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventCard, error)
	SetUserStatusByClubID(clubID int64, userID int64, status string) error
	ApproveRejectUserParticipateInClub(clubID int64, userID int64, decision string) error
	GetUserStatusInClub(clubID int64, userID int64) (*models.ClubUser, error)
}
