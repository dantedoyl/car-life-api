package users_repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/users"
	"github.com/lib/pq"
	"github.com/tarantool/go-tarantool"
	"strconv"
	"sync"
)

type UsersRepository struct {
	sqlConn       *sql.DB
	tarantoolConn *tarantool.Connection
	userSessions  map[string]*models.Session
}

var mtx sync.Mutex

func NewUserRepository(connP *sql.DB, connT *tarantool.Connection) users.IUsersRepository {
	return &UsersRepository{
		sqlConn:       connP,
		tarantoolConn: connT,
		userSessions:  make(map[string]*models.Session),
	}
}

func (ur *UsersRepository) InsertUser(user *models.User, car *models.CarCard) (*models.User, error) {
	_, err := ur.sqlConn.Exec(
		`INSERT INTO users
                (vk_id, name, surname, avatar, tags, description)
                VALUES ($1, $2, $3, $4, $5, $6)`,
		user.VKID,
		user.Name,
		user.Surname,
		user.AvatarUrl,
		pq.Array(user.Tags),
		user.Description)
	if err != nil {
		return nil, err
	}
	if car != nil {
		//if len(user.Garage) != 0 {
		err = ur.sqlConn.QueryRow(
			`INSERT INTO cars
               (owner_id, brand, model,date,description, body, engine, horse_power, name)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
               RETURNING id`,
			user.VKID, car.Brand, car.Model, car.Date, car.Description, car.Body, car.Engine, car.HorsePower, car.Name).Scan(&car.ID)
		if err != nil {
			return nil, err
		}
		//}
		user.CarID = int64(car.ID)
	}
	return user, nil
}

func (ur *UsersRepository) InsertCar(car *models.CarCard) (*models.CarCard, error) {
	err := ur.sqlConn.QueryRow(
		`INSERT INTO cars
               (owner_id, brand, model,date,description, body, engine, horse_power, name)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
               RETURNING id`,
		car.OwnerID, car.Brand, car.Model, car.Date, car.Description, car.Body, car.Engine, car.HorsePower, car.Name).Scan(&car.ID)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (ur *UsersRepository) SelectByID(userID uint64) (*models.User, error) {
	user := &models.User{}
	err := ur.sqlConn.QueryRow(
		`SELECT  vk_id, name, surname, avatar, tags, description from users
				WHERE vk_id = $1`, userID).Scan(&user.VKID, &user.Name, &user.Surname, &user.AvatarUrl, pq.Array(&user.Tags), &user.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	//var cars []*models.CarCard
	//
	//q := `SELECT id, owner_id, brand, model,date,description, avatar, body, engine, horse_power, name FROM cars WHERE owner_id = $1`
	//rows, err := ur.sqlConn.Query(q, userID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//defer rows.Close()
	//
	//for rows.Next() {
	//	car := &models.CarCard{}
	//	err = rows.Scan(&car.ID, &car.OwnerID, &car.Brand, &car.Model, &car.Date, &car.Description, &car.AvatarUrl, user.Garage[0].Body, user.Garage[0].Engine, user.Garage[0].HorsePower, user.Garage[0].Name)
	//	if err != nil {
	//		return nil, err
	//	}
	//	cars = append(cars, car)
	//}
	//user.Garage = cars

	return user, nil
}

func (ur *UsersRepository) Insert(session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	dataStr := string(data)

	_, err = ur.tarantoolConn.Insert("sessions", []interface{}{session.Value, dataStr})
	if err != nil {
		return err
	}
	//______________________________________
	// map session
	//defer mtx.Unlock()
	//mtx.Lock()
	//ur.userSessions[session.Value] = session

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

	//______________________________________
	// map session
	//sess, ok := ur.userSessions[sessValue]
	//if !ok {
	//	return nil, fmt.Errorf("no session")
	//}

	return sess, nil
}

func (ur *UsersRepository) DeleteByValue(sessionValue string) error {
	_, err := ur.tarantoolConn.Delete("sessions", "primary", []interface{}{sessionValue})
	if err != nil {
		return err
	}

	//______________________________________
	// map session
	//delete(ur.userSessions, sessionValue)

	return nil
}

func (ur *UsersRepository) SelectCarByID(carID uint64) (*models.CarCard, error) {
	car := &models.CarCard{}
	err := ur.sqlConn.QueryRow(
		`SELECT id, owner_id, brand, model,date,description, avatar, body, engine, horse_power, name FROM cars
				WHERE id = $1`, carID).Scan(&car.ID, &car.OwnerID, &car.Brand, &car.Model, &car.Date, &car.Description, &car.AvatarUrl, &car.Body, &car.Engine, &car.HorsePower, &car.Name)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (ur *UsersRepository) UpdateCar(car *models.CarCard) (*models.CarCard, error) {
	err := ur.sqlConn.QueryRow(
		`UPDATE cars SET brand = $2, avatar = $3, model = $4, description = $5, date = $6, body = $7, engine = $8, horse_power = $9, name = $10
                WHERE id = $1
                RETURNING id, owner_id, brand, model,date,description, avatar, body, engine, horse_power, name`,
		car.ID,
		car.Brand,
		car.AvatarUrl,
		car.Model,
		car.Description,
		car.Date,
		car.Body, car.Engine, car.HorsePower, car.Name).Scan(&car.ID, &car.OwnerID, &car.Brand, &car.Model, &car.Date, &car.Description, &car.AvatarUrl, &car.Body, &car.Engine, &car.HorsePower, &car.Name)
	if err != nil {
		return nil, err
	}

	return car, nil
}

func (ur *UsersRepository) GetClubsByUserStatus(userID int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.ClubCard, error) {
	var clubs []*models.ClubCard
	ind := 3
	var values []interface{}
	values = append(values, status, userID)
	q := `SELECT c.id, c.name, c.tags, c.participants_count, c.avatar from users_clubs as uc inner join clubs as c on c.id = uc.club_id WHERE uc.status = $1 and uc.user_id = $2 `

	if idGt != nil {
		q += ` AND c.id > $` + strconv.Itoa(ind)
		values = append(values, idGt)
		ind++
	}

	if idLte != nil {
		q += ` AND c.id <= $` + strconv.Itoa(ind)
		values = append(values, idLte)
		ind++
	}

	if limit != nil {
		q += ` LIMIT $` + strconv.Itoa(ind)
		values = append(values, limit)
	}

	q += ` ORDER BY name desc`
	rows, err := ur.sqlConn.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		club := &models.ClubCard{}
		err = rows.Scan(&club.ID, &club.Name, pq.Array(&club.Tags), &club.ParticipantsCount, &club.AvatarUrl)
		if err != nil {
			return nil, err
		}
		clubs = append(clubs, club)
	}
	return clubs, nil
}

func (ur *UsersRepository) SelectCarByUserID(userID int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.CarCard, error) {
	var cars []*models.CarCard
	ind := 2
	var values []interface{}
	values = append(values, userID)
	q := `SELECT c.id, c.owner_id, c.brand, c.model,c.date,c.description, c.avatar, c.body, c.engine, c.horse_power, c.name from cars as c WHERE c.owner_id = $1`

	if idGt != nil {
		q += ` AND c.id > $` + strconv.Itoa(ind)
		values = append(values, idGt)
		ind++
	}

	if idLte != nil {
		q += ` AND c.id <= $` + strconv.Itoa(ind)
		values = append(values, idLte)
		ind++
	}

	if limit != nil {
		q += ` LIMIT $` + strconv.Itoa(ind)
		values = append(values, limit)
	}

	q += ` ORDER BY c.name desc`
	rows, err := ur.sqlConn.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		car := &models.CarCard{}
		err = rows.Scan(&car.ID, &car.OwnerID, &car.Brand, &car.Model, &car.Date, &car.Description, &car.AvatarUrl, &car.Body, &car.Engine, &car.HorsePower, &car.Name)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

func (ur *UsersRepository) GetEventsByUserStatus(userID int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventCard, error) {
	var events []*models.EventCard
	ind := 3
	var values []interface{}
	values = append(values, status, userID)
	q := `SELECT e.id, e.name, e.event_date, e.latitude, e.longitude, e.avatar from users_events as ue inner join events as e on e.id = ue.event_id WHERE ue.status = $1 and ue.user_id = $2`

	if idGt != nil {
		q += ` AND e.id > $` + strconv.Itoa(ind)
		values = append(values, idGt)
		ind++
	}

	if idLte != nil {
		q += ` AND e.id <= $` + strconv.Itoa(ind)
		values = append(values, idLte)
		ind++
	}

	if limit != nil {
		q += ` LIMIT $` + strconv.Itoa(ind)
		values = append(values, limit)
	}

	q += ` ORDER BY e.event_date desc`
	rows, err := ur.sqlConn.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		event := &models.EventCard{}
		err = rows.Scan(&event.ID, &event.Name, &event.EventDate,
			&event.Latitude, &event.Longitude, &event.AvatarUrl)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
