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

	return nil
}

func (cr *ClubsRepository) GetClubByID(id int64) (*models.Club, error) {
	club := &models.Club{}
	err := cr.dbConn.QueryRow(
		`SELECT  id, name, description, tags, events_count, participants_count, avatar from clubs
				WHERE id = $1`, id).Scan(&club.ID, &club.Name, pq.Array(&club.Tags), &club.EventsCount, &club.ParticipantsCount, &club.AvatarUrl)
	if err != nil {
		return nil, err
	}
	return club, nil
}

func (cr *ClubsRepository) GetClubs(idGt  *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Club, error) {
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
		q += ` AND (name like '%' || $`+strconv.Itoa(ind)+` || '%' OR '%' || $`+strconv.Itoa(ind)+` || '%' like any(tags)`
		values = append(values, idLte)
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
		err = rows.Scan(&club.ID, &club.Name, pq.Array(&club.Tags), &club.EventsCount, &club.ParticipantsCount, &club.AvatarUrl)
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

func (cr *ClubsRepository) GetTags() ([]models.Tag, error) {
	var tags []models.Tag
	rows, err := cr.dbConn.Query(`SELECT name from tags ORDER BY usage_count desc`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

