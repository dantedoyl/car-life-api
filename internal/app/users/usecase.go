package users

import (
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type IUsersUsecase interface {
	CreateSession(sess *models.Session) error
	GetSession(sessValue string) (*models.Session, error)
	DeleteSession(sessionValue string) error
	CheckSession(sessValue string) (*models.Session, error)

	Create(user *models.User, car *models.CarCard) (*models.User, error)
	GetByID(vkID uint64) (*models.User, error)
	GetClubsByUserStatus(userID int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.ClubCard, error)
	UpdateAvatar(carID uint64, fileHeader *multipart.FileHeader) (*models.User, error)
	SelectCarByUserID(userID int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.CarCard, error)
	GetEventsByUserStatus(userID int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventCard, error)
}
