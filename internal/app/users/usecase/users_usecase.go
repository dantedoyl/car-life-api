package usecase

import (
	"github.com/dantedoyl/car-life-api/internal/app/clients/filesystem"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/users"
	"mime/multipart"
	"time"
)

type UsersUsecase struct {
	usersRepo users.IUsersRepository
}

func NewUsersUsecase(repo users.IUsersRepository) users.IUsersUsecase {
	return &UsersUsecase{
		usersRepo: repo,
	}
}

func (uu *UsersUsecase) Create(user *models.User) error {
	err := uu.usersRepo.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UsersUsecase) GetByID(vkID uint64) (*models.User, error) {
	user, err := uu.usersRepo.SelectByID(vkID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (uu *UsersUsecase) CreateSession(sess *models.Session) error {
	err := uu.usersRepo.Insert(sess)
	if err != nil {
		return err
	}

	return nil
}

func (uu *UsersUsecase) GetSession(sessValue string) (*models.Session, error) {
	sess, err := uu.usersRepo.SelectByValue(sessValue)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (uu *UsersUsecase) DeleteSession(sessionValue string) error {
	if _, err := uu.GetSession(sessionValue); err != nil {
		return err
	}

	err := uu.usersRepo.DeleteByValue(sessionValue)
	if err != nil {
		return err
	}

	return nil
}

func (uu *UsersUsecase) CheckSession(sessValue string) (*models.Session, error) {
	sess, err := uu.GetSession(sessValue)
	if err != nil {
		return nil, err
	}

	if sess.ExpiresAt.Before(time.Now()) {
		err := uu.DeleteSession(sessValue)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	return sess, nil
}

func (uu *UsersUsecase) UpdateAvatar(carID uint64, fileHeader *multipart.FileHeader) (*models.User, error) {
	car, err := uu.usersRepo.SelectCarByID(carID)
	if err != nil {
		return nil, err
	}

	imgUrl, err := filesystem.InsertPhoto(fileHeader, "img/cars/")
	if err != nil {
		return nil, err
	}

	oldAvatar := car.AvatarUrl
	car.AvatarUrl = imgUrl
	car, err = uu.usersRepo.UpdateCar(car)
	if err != nil {
		return nil, err
	}

	user, err := uu.usersRepo.SelectByID(car.OwnerID)
	if err != nil {
		return nil, err
	}

	if oldAvatar == "/img/cars/default.jpeg" {
		return user, nil
	}

	err = filesystem.RemovePhoto(oldAvatar)
	if err != nil {
		return nil, err
	}

	return user, nil
}

