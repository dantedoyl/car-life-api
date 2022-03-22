package events_repository

import (
	"database/sql"
	"github.com/dantedoyl/car-life-api/internal/app/events"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"strconv"
)

type EventsRepository struct {
	dbConn *sql.DB
}

func NewProductRepository(conn *sql.DB) events.IEventsRepository {
	return &EventsRepository{
		dbConn: conn,
	}
}

func (er *EventsRepository) InsertEvent(event *models.Event) error {
	err := er.dbConn.QueryRow(
		`INSERT INTO events
                (name, club_id, description, event_date, latitude, longitude)
                VALUES ($1, $2, $3, $4, $5, $6) 
                RETURNING id`,
		event.Name,
		event.Club.ID,
		event.Description,
		event.EventDate,
		event.Latitude,
		event.Longitude).Scan(&event.ID)
	if err != nil {
		return err
	}

	return nil
}

func (er *EventsRepository) GetEventByID(id int64) (*models.Event, error) {
	event := &models.Event{}
	err := er.dbConn.QueryRow(
		`SELECT  id, name, club_id, description, event_date, latitude, longitude, avatar from events
				WHERE id = $1`, id).Scan(&event.ID, &event.Name, &event.Club.ID, &event.Description, &event.EventDate,
		&event.Latitude, &event.Longitude, &event.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (er *EventsRepository) GetEvents(idGt  *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Event, error) {
	var events []*models.Event
	ind := 1
	var values []interface{}
	q := `SELECT  id, name, club_id, description, event_date, latitude, longitude, avatar from events WHERE true `

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
		q += ` AND name like '%' || $`+strconv.Itoa(ind)+` || '%'`
		values = append(values, idLte)
		ind++
	}

	if limit != nil {
		q += ` LIMIT $` + strconv.Itoa(ind)
		values = append(values, limit)
	}

	q += ` ORDER BY created_at desc`
	rows, err := er.dbConn.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		event := &models.Event{}
		err = rows.Scan(&event.ID, &event.Name, &event.Club.ID, &event.Description, &event.EventDate,
			&event.Latitude, &event.Longitude, &event.AvatarUrl)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (er *EventsRepository) UpdateEvent(event *models.Event) (*models.Event, error) {
	err := er.dbConn.QueryRow(
		`UPDATE events SET name = $1, description = $2, event_date = $3, latitude = $4, longitude = $5, avatar = $6
				WHERE id = $7
				RETURNING id, name, club_id, description, event_date, latitude, longitude, avatar`,
		event.Name, event.Description, event.EventDate, event.Latitude, event.Longitude, event.AvatarUrl, event.ID).Scan(&event.ID, &event.Name, &event.Club.ID, &event.Description, &event.EventDate,
		&event.Latitude, &event.Longitude, &event.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return event, nil
}
