package events_repository

import (
	"database/sql"
	"github.com/dantedoyl/car-life-api/internal/app/events"
	"github.com/dantedoyl/car-life-api/internal/app/models"
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
                VALUES ($1, $2, $3, $4, $5) 
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
		`SELECT  id, name, club_id, description, event_date, latitude, longitude from events
				WHERE id = $1`, id).Scan(&event.ID, &event.Name, &event.Club.ID, &event.Description, &event.EventDate,
					&event.Latitude, &event.Longitude)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (er *EventsRepository) GetEvents() ([]*models.Event, error) {
	var events []*models.Event
	rows, err := er.dbConn.Query(`SELECT  id, name, club_id, description, event_date, latitude, longitude from events`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		event := &models.Event{}
		err = rows.Scan(&event.ID, &event.Name, &event.Club.ID, &event.Description, &event.EventDate,
			&event.Latitude, &event.Longitude)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}