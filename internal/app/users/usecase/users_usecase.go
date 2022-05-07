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

func (uu *UsersUsecase) Create(user *models.User, car *models.CarCard) (*models.User, error) {
	user, err := uu.usersRepo.InsertUser(user, car)
	if err != nil {
		return nil, err
	}
	return user, nil
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

func (uu *UsersUsecase) GetClubsByUserStatus(userID int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.ClubCard, error) {
	return uu.usersRepo.GetClubsByUserStatus(userID, status, idGt, idLte, limit)
}

func (uu *UsersUsecase) SelectCarByUserID(userID int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.CarCard, error){
return uu.usersRepo.SelectCarByUserID(userID, idGt, idLte, limit)
}

func (uu *UsersUsecase) GetEventsByUserStatus(userID int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventCard, error) {
	return uu.usersRepo.GetEventsByUserStatus(userID, status, idGt, idLte, limit)
}

func (uu *UsersUsecase) AddNewUserCar(car *models.CarCard) (*models.CarCard, error) {
	return uu.usersRepo.InsertCar(car)
}

func (uu *UsersUsecase) UpdateUserInfo(user *models.User) (*models.User, error) {
	return uu.usersRepo.Update(user)
}

func(uu *UsersUsecase)	SelectCarByID(carID int64) (*models.CarCard, error) {
	return uu.usersRepo.SelectCarByID(uint64(carID))
}

func(uu *UsersUsecase)	DeleteCarByID(carID int64) error {
	return uu.usersRepo.DeleteCarByID(carID)
}

func(uu *UsersUsecase)	ComplainByID(target string, complaint models.Complaint) error {
	return uu.usersRepo.ComplainByID(target, complaint)
}
