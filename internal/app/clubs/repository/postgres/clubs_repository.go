package clubs_repository

import (
	"database/sql"
	clubs "github.com/dantedoyl/car-life-api/internal/app/clubs"
	"github.com/dantedoyl/car-life-api/internal/app/models"
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
                (name, description)
                VALUES ($1, $2) 
                RETURNING id`,
		club.Name,
		club.Description).Scan(&club.ID)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ClubsRepository) GetClubByID(id int64) (*models.Club, error) {
	club := &models.Club{}
	err := cr.dbConn.QueryRow(
		`SELECT  id, name, description, events_count, participants_count, avatar from clubs
				WHERE id = $1`, id).Scan(&club.ID, &club.Name, &club.EventsCount, &club.ParticipantsCount, &club.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return club, nil
}

func (cr *ClubsRepository) GetClubs() ([]*models.Club, error) {
	var clubs []*models.Club
	rows, err := cr.dbConn.Query(`SELECT  id, name, description, events_count, participants_count, avatar from clubs ORDER BY created_at desc`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		club := &models.Club{}
		err = rows.Scan(&club.ID, &club.Name, &club.EventsCount, &club.ParticipantsCount, &club.AvatarUrl)
		if err != nil {
			return nil, err
		}
		clubs = append(clubs, club)
	}
	return clubs, nil
}

func (cr *ClubsRepository) UpdateClub(club *models.Club) (*models.Club, error) {
	err := cr.dbConn.QueryRow(
		`UPDATE clubs SET name = $1, description = $2, avatar = $3 from events
				WHERE id = $1
				RETURNING id, name, description, events_count, participants_count, avatar`,
				club.Name, club.Description, club.AvatarUrl).Scan(&club.ID, &club.Name, &club.EventsCount, &club.ParticipantsCount, &club.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return club, nil
}

