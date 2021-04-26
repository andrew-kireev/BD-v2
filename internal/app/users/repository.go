package users

import (
	"BD-v2/internal/app/users/models"
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	CheckIfUserExist(ctx context.Context, user *models.User) ([]*models.User, error)
}
