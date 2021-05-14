package repository

import (
	"BD-v2/internal/app/forums/models"
	models2 "BD-v2/internal/app/threads/models"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type ForumRepository struct {
	db *sqlx.DB
}

func NewForumRepository(db *sqlx.DB) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (rep *ForumRepository) CreateForum(forum *models.Forum) error {
	query := `insert into forums (title, users, slug) values ($1, $2, $3)`

	_, err := rep.db.Exec(query, forum.Title, forum.User, forum.Slug)
	return err
}

func (rep *ForumRepository) GetForumSlug(slug string) (*models.Forum, error) {
	query := `select title, users, slug, posts, threads from forums
		where slug = $1`
	forum := &models.Forum{}
	err := rep.db.Get(forum, query, slug)

	return forum, err
}

func (rep *ForumRepository) ClearDB() error {
	query := `TRUNCATE TABLE users CASCADE;
		TRUNCATE TABLE forums CASCADE;
		TRUNCATE TABLE threads CASCADE;`

	_, err := rep.db.Exec(query)
	return err
}

func (rep *ForumRepository) GetTreads(limit int, forum, since string, desc bool) ([]*models2.Thread, error) {
	query := `select id, title, author, forum, message, votes, slug, created from threads
		where forum = $1
		`
	if since != "" && desc {
		query += fmt.Sprintf(" and created <= '%s'", since)
	} else if since != "" {
		query += fmt.Sprintf(" and created >= '%s'", since)
	}

	if desc {
		query += " order by created desc"
	} else {
		query += " order by created"
	}

	query += fmt.Sprintf(" limit NULLIF(%d, 0)", limit)

	threads := make([]*models2.Thread, 0)
	rows, err := rep.db.Query(query, forum)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := &time.Time{}
		thread := &models2.Thread{}
		err = rows.Scan(&thread.ID, &thread.Title,
			&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
			&thread.Slug, t)
		thread.Created = strfmt.DateTime(t.UTC()).String()
		threads = append(threads, thread)
	}
	return threads, err
}

