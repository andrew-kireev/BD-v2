package threads

import (
	models2 "BD-v2/internal/app/posts/models"
	"BD-v2/internal/app/threads/models"
	"context"
)

type Repository interface {
	CreateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error)
	FindThreadSlug(ctx context.Context, slug string) (*models.Thread, error)
	FindThreadID(ctx context.Context, threadID int) (*models.Thread, error)
	CreatePost(ctx context.Context, post *models2.Post) (*models2.Post, error)
	AddVoice(ctx context.Context, voice *models.Vote) error
	GetVoice(ctx context.Context, voice *models.Vote) (*models.Vote, error)
	UpdateVoice(ctx context.Context, voice *models.Vote) error
}
