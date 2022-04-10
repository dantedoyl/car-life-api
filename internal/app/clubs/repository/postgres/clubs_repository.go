package clubs_repository

import (
	"database/sql"
	clubs "github.com/dantedoyl/car-life-api/internal/app/clubs"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/lib/pq"
	"strconv"
)

type ClubsRepository struct {
	dbConn *sql.DB
}

func NewClubRepository(conn *sql.DB) clubs.IClubsRepository {
	return &ClubsRepository{
		dbConn: conn,
	}
}

func (cr *ClubsRepository) InsertClub(club *models.Club) error {
	err := cr.dbConn.QueryRow(
		`INSERT INTO clubs
                (name, description, tags)
                VALUES ($1, $2, $3) 
                RETURNING id`,
		club.Name,
		club.Description,
		pq.Array(club.Tags)).Scan(&club.ID)
	if err != nil {
		return err
	}

	_, err = cr.dbConn.Exec(
		`UPDATE tags SET usage_count = usage_count + 1 WHERE name = any($1)`,
		pq.Array(club.Tags))
	if err != nil {
		return err
	}

	_, err = cr.dbConn.Exec(
		`INSERT INTO users_clubs (club_id, user_id, status) VALUES ($1, $2, $3)
				ON CONFLICT (user_id, club_id) DO UPDATE
			SET status = $3`, club.ID, club.OwnerID, "admin")
	if err != nil {
		return err
	}

	return nil
}

func (cr *ClubsRepository) GetClubByID(id int64, userID uint64) (*models.Club, error) {
	club := &models.Club{}
	err := cr.dbConn.QueryRow(
		`SELECT  c.id, c.name, c.description, c.tags, c.events_count, c.participants_count, c.avatar, uc.user_id as owner_id from clubs as c inner join users_clubs as uc on uc.club_id = c.id
				WHERE c.id = $1 and uc.user_id = 'admin'`, id).Scan(&club.ID, &club.Name, &club.Description, pq.Array(&club.Tags), &club.EventsCount, &club.ParticipantsCount, &club.AvatarUrl, &club.OwnerID)
	if err != nil {
		return nil, err
	}

	if userID != 0 {
		var status string
		err = cr.dbConn.QueryRow(
			`SELECT status from users_clubs
				WHERE club_id = $1 and user_id = $2`, club.ID, userID).Scan(&status)
		if err == sql.ErrNoRows {
			club.UserStatus = "unknown"
			return club, nil
		}
		if err != nil {
			return nil, err
		}
		club.UserStatus = status
	}

	return club, nil
}

func (cr *ClubsRepository) GetClubs(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Club, error) {
	var clubs []*models.Club
	ind := 1
	var values []interface{}
	q := `SELECT  id, name, description, tags, events_count, participants_count, avatar from clubs WHERE true `

	if idGt != nil {
		q += ` AND id > $` + strconv.Itoa(ind)
		values = append(values, idGt)
		ind++
	}

	if idLte != nil {
		q += ` AND id <= $` + strconv.Itoa(ind)
		values = append(values, idLte)
		ind++
	}

	if query != nil {
		q += ` AND (name like '%' || $` + strconv.Itoa(ind) + ` || '%' OR EXISTS (
    			SELECT 
				FROM   unnest(tags) elem
   				 WHERE  elem LIKE '%' || $` + strconv.Itoa(ind) + ` || '%'))`
		values = append(values, query)
		ind++
	}

	if limit != nil {
		q += ` LIMIT $` + strconv.Itoa(ind)
		values = append(values, limit)
	}

	q += ` ORDER BY created_at desc`
	rows, err := cr.dbConn.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		club := &models.Club{}
		err = rows.Scan(&club.ID, &club.Name, &club.Description, pq.Array(&club.Tags), &club.EventsCount, &club.ParticipantsCount, &club.AvatarUrl)
		if err != nil {
			return nil, err
		}
		clubs = append(clubs, club)
	}
	return clubs, nil
}

func (cr *ClubsRepository) UpdateClub(club *models.Club) (*models.Club, error) {
	err := cr.dbConn.QueryRow(
		`UPDATE clubs SET name = $1, description = $2, avatar = $3
				WHERE id = $4
				RETURNING id, name, description, events_count, participants_count, avatar`,
		club.Name, club.Description, club.AvatarUrl, club.ID).Scan(&club.ID, &club.Name, &club.Description, &club.EventsCount, &club.ParticipantsCount, &club.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return club, nil
}

func (cr *ClubsRepository) GetTags() ([]models.Tag, error) {
	var tags []models.Tag
	rows, err := cr.dbConn.Query(`SELECT id, name from tags ORDER BY usage_count desc`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (cr *ClubsRepository) GetClubsUserByStatus(club_id int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.UserCard, error) {
	var users []*models.UserCard
	ind := 3
	var values []interface{}
	values = append(values, status, club_id)
	q := `SELECT u.vk_id, u.name, u.surname, u.avatar from users_clubs as uc INNER JOIN users as u on u.vk_id = uc.user_id WHERE uc.status = $1 and uc.club_id=$2`

	if idGt != nil {
		q += ` AND u.vk_id > $` + strconv.Itoa(ind)
		values = append(values, idGt)
		ind++
	}

	if idLte != nil {
		q += ` AND u.vk_id <= $` + strconv.Itoa(ind)
		values = append(values, idLte)
		ind++
	}

	if limit != nil {
		q += ` LIMIT $` + strconv.Itoa(ind)
		values = append(values, limit)
	}

	q += ` ORDER BY u.surname desc`
	rows, err := cr.dbConn.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := &models.UserCard{}
		err = rows.Scan(&user.VKID, &user.Name, &user.Surname, &user.AvatarUrl)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (cr *ClubsRepository) GetClubsCars(club_id int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.CarCard, error) {
	var cars []*models.CarCard
	ind := 2
	var values []interface{}
	values = append(values, club_id)
	q := `SELECT c.id, c.owner_id, c.brand, c.model,c.date,c.description, c.avatar, c.body, c.engine, c.horse_power, c.name from cars as c JOIN users_clubs as uc on c.owner_id = uc.user_id WHERE uc.status = 'participant' and uc.club_id=$1`

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
	rows, err := cr.dbConn.Query(q, values...)
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

func (cr *ClubsRepository) GetClubsEvents(club_id int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventCard, error) {
	var events []*models.EventCard
	ind := 2
	var values []interface{}
	values = append(values, club_id)
	q := `SELECT e.id, e.name, e.event_date, e.latitude, e.longitude, e.avatar from events as e WHERE e.club_id=$1`

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
	rows, err := cr.dbConn.Query(q, values...)
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

func (cr *ClubsRepository) SetUserStatusByClubID(clubID int64, userID int64, status string) error {
	_, err := cr.dbConn.Exec(
		`INSERT INTO users_clubs (club_id, user_id, status) VALUES ($1, $2, $3)
				ON CONFLICT (user_id, club_id) DO UPDATE
			SET status = $3`, clubID, userID, status)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ClubsRepository) GetUserStatusInClub(clubID int64, userID int64) (*models.ClubUser, error) {
	userClub := &models.ClubUser{}
	err := cr.dbConn.QueryRow(`SELECT club_id, user_id, status FROM users_clubs WHERE club_id = $1 and user_id = $2`, clubID, userID).Scan(
		&userClub.ClubID, &userClub.UserID, &userClub.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return userClub, nil
}

func (cr *ClubsRepository) GetClubChatID(clubID int64, userID int64) (int64, error) {
	var chatID int64
	err := cr.dbConn.QueryRow(`SELECT c.chat_id FROM clubs as c inner join users_clubs as uc on c.id = uc.club_id WHERE c.id = $1 and uc.user_id = $2`, clubID, userID).Scan(&chatID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

func (cr *ClubsRepository) SetClubChatID(clubID int64, chatID int64) error {
	_, err := cr.dbConn.Exec(`UPDATE clubs SET chat_id = $1 WHERE id = $2`, chatID, clubID)
	if err != nil {
		return err
	}
	return nil
}
