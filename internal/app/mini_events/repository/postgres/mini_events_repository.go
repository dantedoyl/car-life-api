package mini_events_repository

import (
	"database/sql"
	"github.com/dantedoyl/car-life-api/internal/app/mini_events"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"strconv"
)

type MiniEventsRepository struct {
	dbConn *sql.DB
}

func NewMiniEventsRepository(conn *sql.DB) mini_events.IMiniEventsRepository {
	return &MiniEventsRepository{
		dbConn: conn,
	}
}

func (mr *MiniEventsRepository) InsertMiniEvent(event *models.MiniEvent) error {
	err := mr.dbConn.QueryRow(
		`INSERT INTO events
                (type_id, user_id, description, created_at, ended_at, latitude, longitude)
                VALUES ($1, $2, $3, $4, $5, $6, $7) 
                RETURNING id`,
		event.Type.ID,
		event.User.VKID,
		event.Description,
		event.CreatedAt,
		event.EndedAt,
		event.Latitude,
		event.Longitude).Scan(&event.ID)
	if err != nil {
		return err
	}

	return nil
}

func (mr *MiniEventsRepository) GetMiniEventByID(id int64) (*models.MiniEvent, error) {
	event := &models.MiniEvent{}
	err := mr.dbConn.QueryRow(
		`SELECT  id, type_id, user_id, description, created_at, ended_at, latitude, longitude, from mini_events
				WHERE id = $1`, id).Scan(&event.ID, &event.Type.ID, &event.User.VKID, &event.Description, &event.CreatedAt, &event.EndedAt,
		&event.Latitude, &event.Longitude)
	if err != nil {
		return nil, err
	}

	err = mr.dbConn.QueryRow(
		`SELECT vk_id, name, surname, avatar from users
				WHERE vk_id = $1`, event.User.VKID).Scan(&event.User.VKID, &event.User.Name, &event.User.Surname, &event.User.AvatarUrl)
	if err != nil {
		return nil, err
	}

	err = mr.dbConn.QueryRow(
		`SELECT id, public_name, public_description from mini_event_type
				WHERE id = $1`, event.Type.ID).Scan(&event.Type.ID, &event.Type.PublicName, &event.Type.PublicDescription)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (mr *MiniEventsRepository) GetMiniEvents(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.MiniEvent, error) {
	var events []*models.MiniEvent
	ind := 1
	var values []interface{}
	q := `SELECT me.id, me.type_id, met.public_name, met.public_description, me.user_id, u.name, u.surname, u.avatar, me.description, me.created_at, me.ended_at, me.latitude, me.longitude, from mini_events as me 
			left join mini_event_type as met on me.type_id = met.id 
			left join users as u me.user_id = u.vk_id
			WHERE true `

	if idGt != nil {
		q += ` AND me.id > $` + strconv.Itoa(ind)
		values = append(values, idGt)
		ind++
	}

	if idLte != nil {
		q += ` AND me.id <= $` + strconv.Itoa(ind)
		values = append(values, idLte)
		ind++
	}

	if limit != nil {
		q += ` LIMIT $` + strconv.Itoa(ind)
		values = append(values, limit)
	}

	q += ` ORDER BY created_at desc`
	rows, err := mr.dbConn.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		event := &models.MiniEvent{}
		err = rows.Scan(&event.ID, &event.Type.ID, &event.Type.PublicName, &event.Type.PublicDescription,
			&event.User.VKID, &event.User.Name, &event.User.Surname, &event.User.AvatarUrl,
			&event.Description, &event.CreatedAt, &event.EndedAt,
			&event.Latitude, &event.Longitude)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
