package users_repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/users"
	"github.com/tarantool/go-tarantool"
)

type UsersRepository struct {
	sqlConn *sql.DB
	tarantoolConn *tarantool.Connection
}

func NewUserRepository(connP *sql.DB, connT *tarantool.Connection) users.IUsersRepository {
	return &UsersRepository{
		sqlConn: connP,
		tarantoolConn: connT,
	}
}

func (ur *UsersRepository) InsertUser(user *models.User) error {
	_, err := ur.sqlConn.Exec(
		`INSERT INTO users
                (vk_id, name, surname, avatar)
                VALUES ($1, $2, $3, $4)`,
		user.VKID,
		user.Name,
		user.Surname,
		user.AvatarUrl)
	if err != nil {
		return err
	}

	err = ur.sqlConn.QueryRow(
		`INSERT INTO cars
                (name, owner_id)
                VALUES ($1, $2)`,
		user.Garage[0].Name,
		user.VKID).Scan(&user.Garage[0].ID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UsersRepository) SelectByID(userID uint64) (*models.User, error) {
	user := &models.User{}
	err := ur.sqlConn.QueryRow(
		`SELECT  vk_id, name, surname, avatar from users
				WHERE vk_id = $1`, userID).Scan(&user.VKID, &user.Name, &user.Surname, &user.AvatarUrl)
	if err != nil {
		return nil, err
	}

	var cars []*models.CarCard

	q := `SELECT id, name, owner_id, avatar FROM cars WHERE owner_id = $1`
	rows, err := ur.sqlConn.Query(q, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		car := &models.CarCard{}
		err = rows.Scan(&car.ID, &car.Name, &car.OwnerID, &car.AvatarUrl)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	user.Garage = cars

	return user, nil
}

func (ur *UsersRepository) Insert(session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	dataStr := string(data)

	//resp, err := sr.dbConn.Eval("return new_session(...)", []interface{}{session.Value, dataStr})
	_, err = ur.tarantoolConn.Insert("sessions", []interface{}{session.Value, dataStr})
	if err != nil {
		return err
	}

	return nil
}

func (ur *UsersRepository) SelectByValue(sessValue string) (*models.Session, error) {
	resp, err := ur.tarantoolConn.Call("check_session", []interface{}{sessValue})
	if err != nil {
		return nil, err
	}

	data := resp.Data[0]
	if data == nil {
		return &models.Session{}, nil
	}

	sessionDataSlice, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot cast data")
	}

	if sessionDataSlice[0] == nil {
		return nil, fmt.Errorf("session not exist")
	}

	sessionData, ok := sessionDataSlice[1].(string)
	if !ok {
		return nil, fmt.Errorf("cannot cast to string")
	}

	sess := &models.Session{}
	err = json.Unmarshal([]byte(sessionData), sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (ur *UsersRepository) DeleteByValue(sessionValue string) error {
	_, err := ur.tarantoolConn.Delete("sessions", "primary", []interface{}{sessionValue})
	if err != nil {
		return err
	}

	return nil
}

func (ur *UsersRepository)SelectCarByID(carID uint64) (*models.CarCard, error) {
	car := &models.CarCard{}
	err := ur.sqlConn.QueryRow(
		`SELECT id, owner_id, name, avatar from cars
				WHERE id = $1`, carID).Scan(&car.ID, &car.OwnerID, &car.Name, &car.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (ur *UsersRepository) UpdateCar(car *models.CarCard) (*models.CarCard, error) {
	err := ur.sqlConn.QueryRow(
		`UPDATE cars SET name = $2, avatar = $3
                WHERE id = $1
                RETURNING id, name, avatar, owner_id)`,
		car.ID,
		car.Name,
		car.AvatarUrl).Scan(&car.ID, &car.Name, &car.AvatarUrl, &car.OwnerID)
	if err != nil {
		return nil, err
	}

	return car, nil
}
