package users

import "github.com/dantedoyl/car-life-api/internal/app/models"

type IUsersRepository interface {
	Insert(session *models.Session) error
	SelectByValue(sessValue string) (*models.Session, error)
	DeleteByValue(sessionValue string) error

	InsertUser(user *models.User, car *models.CarCard) (*models.User, error)
	SelectByID(userID uint64) (*models.User, error)
	SelectCarByID(carID uint64) (*models.CarCard, error)
	UpdateCar(car *models.CarCard) (*models.CarCard, error)
	//Update(user *models.User) error
	//Delete(userID uint64) error
}
