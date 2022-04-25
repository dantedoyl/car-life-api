package events_posts_repository

import (
	"database/sql"
	"fmt"
	"github.com/dantedoyl/car-life-api/internal/app/events_posts"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/lib/pq"
	"strconv"
)

type EventsPostsRepository struct {
	dbConn *sql.DB
}

func NewEventsPostsRepository(conn *sql.DB) events_posts.IEventsPostsRepository {
	return &EventsPostsRepository{
		dbConn: conn,
	}
}

func (epr *EventsPostsRepository) InsertEventPost(eventPost *models.EventPost) error {
	err := epr.dbConn.QueryRow(
		`INSERT INTO events_posts
                (text, user_id, event_id)
                VALUES ($1, $2, $3) 
                RETURNING id, created_at`,
		eventPost.Text,
		eventPost.User.VKID,
		eventPost.EventID).Scan(&eventPost.ID, &eventPost.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (epr *EventsPostsRepository) GetEventsPostsByEventID(eventID uint64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventPost, error) {
	var eventsPosts []*models.EventPost
	ind := 1
	var values []interface{}
	q := `SELECT ep.id, ep.text, ep.user_id, u.name, u.surname, u.avatar, ep.event_id, ep.created_at, array_agg(epa.url) from events_posts as ep 
    		left join events_posts_attachments as epa on ep.id = epa.post_id
			left join users as u on u.vk_id = ep.user_id
			WHERE true `

	if idGt != nil {
		q += ` AND ep.id > $` + strconv.Itoa(ind)
		values = append(values, idGt)
		ind++
	}

	if idLte != nil {
		q += ` AND ep.id <= $` + strconv.Itoa(ind)
		values = append(values, idLte)
		ind++
	}

	if limit != nil {
		q += ` LIMIT $` + strconv.Itoa(ind)
		values = append(values, limit)
	}

	q += ` GROUP BY ep.id, u.name ORDER BY ep.created_at desc`
	rows, err := epr.dbConn.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		post := &models.EventPost{}
		err = rows.Scan(&post.ID, &post.Text, &post.User.VKID, &post.User.Name, &post.User.Surname, &post.User.VKID, &post.EventID, &post.CreatedAt, pq.Array(&post.Attachments))
		if err != nil {
			return nil, err
		}
		eventsPosts = append(eventsPosts, post)
	}
	return eventsPosts, nil
}

func (epr *EventsPostsRepository) InsertEventPostAttachments(postID uint64, attachments []string) error {
	ind := 1
	var values []interface{}
	query := `INSERT INTO events_posts_attachments (url, post_id) VALUES`

	for i, attachment := range attachments {
		if i > 0 {
			query += `,`
		}
		query += fmt.Sprintf(` ($%d, $%d)`, ind, ind+1)
		values = append(values, attachment, postID)
		ind = ind + 2
	}

	_, err := epr.dbConn.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (epr *EventsPostsRepository) GetEventPostByPostID(postID uint64) (*models.EventPost, error) {
	post := &models.EventPost{}
	err := epr.dbConn.QueryRow(
		`SELECT ep.id, ep.text, ep.user_id, u.name, u.surname, u.avatar, ep.event_id, ep.created_at, array_agg(epa.url) from events_posts as ep 
    left join events_posts_attachments as epa on ep.id = epa.post_id
			left join users as u on u.vk_id = ep.user_id
			WHERE ep.id = $1 `, postID).Scan(&post.ID, &post.Text, &post.User.VKID, &post.User.Name, &post.User.Surname, &post.User.VKID, &post.EventID, &post.CreatedAt, pq.Array(&post.Attachments))
	if err != nil {
		return nil, err
	}

	return post, nil
}
