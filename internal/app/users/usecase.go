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

	Create(user *models.User) error
	GetByID(vkID uint64) (*models.User, error)

	UpdateAvatar(carID uint64, fileHeader *multipart.FileHeader) (*models.User, error)
}
