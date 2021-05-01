package threads

import (
	"BD-v2/internal/app/threads/models"
	"context"
)

type Repository interface {
	CreateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error)
	FindThreadSlug(ctx context.Context, slug string) (*models.Thread, error)
}
