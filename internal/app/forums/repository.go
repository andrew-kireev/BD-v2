package forums

import "BD-v2/internal/app/forums/models"

type Repository interface {
	CreateForum(forum *models.Forum) error
	GetForumSlug(slug string) (*models.Forum, error)
}
