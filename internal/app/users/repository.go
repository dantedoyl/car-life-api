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
	GetClubsByUserStatus(userID int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.ClubCard, error)
	SelectCarByUserID(userID int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.CarCard, error)
	GetEventsByUserStatus(userID int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventCard, error)
	InsertCar(car *models.CarCard) (*models.CarCard, error)
	Update(user *models.User) (*models.User, error)
	DeleteCarByID(carID int64) error
	ComplainByID(target string, complaint models.Complaint) error
}
