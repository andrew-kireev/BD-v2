package repository

import (
	"BD-v2/internal/app/forums/models"
	"github.com/jmoiron/sqlx"
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
