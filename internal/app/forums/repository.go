package forums

import (
	"BD-v2/internal/app/forums/models"
	models2 "BD-v2/internal/app/threads/models"
)

type Repository interface {
	CreateForum(forum *models.Forum) error
	GetForumSlug(slug string) (*models.Forum, error)
	ClearDB() error
	GetTreads(limit int, forum, since string,  desc bool) ([]*models2.Thread, error)
}
