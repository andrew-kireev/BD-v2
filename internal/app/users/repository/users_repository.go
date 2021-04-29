package repository

import (
	"BD-v2/internal/app/users/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UsersRepository struct {
	DBPool *pgxpool.Pool
}

func NewUsersRepository(db *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{
		DBPool: db,
	}
}

func (rep *UsersRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `insert into users (nickname, fullname, about, email)
	VALUES ($1, $2, $3, $4)`

	_, err := rep.DBPool.Exec(ctx, query, user.Nickname, user.FullName,
		user.About, user.Email)
	return err
}

func (rep *UsersRepository) CheckIfUserExist(ctx context.Context,
	user *models.User) ([]*models.User, error) {
	query := `select nickname, fullname, about, email from users
	where nickname = $1 or email = $2`

	rows, err := rep.DBPool.Query(ctx, query, user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	existedUsers := make([]*models.User, 0)

	for rows.Next() {
		user = &models.User{}
		_ = rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)

		existedUsers = append(existedUsers, user)
	}
	return existedUsers, nil
}

func (rep *UsersRepository) FindUserNickname(ctx context.Context, nickname string) (*models.User, error) {
	query := `select nickname, fullname, about, email from users
		where nickname = $1`
	user := &models.User{}

	err := rep.DBPool.QueryRow(ctx, query, nickname).Scan(&user.Nickname, &user.FullName,
		&user.About, &user.Email)
	return user, err
}

func (rep *UsersRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `update users set fullname = COALESCE(NULLIF($1, ''), fullname),
		email = COALESCE(NULLIF($2, ''), email),
		about = COALESCE(NULLIF($3, ''), about)
		where nickname = $4 returning nickname, fullname, about, email`

	err := rep.DBPool.QueryRow(ctx, query, user.FullName, user.Email,
		user.About, user.Nickname).Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)
	return user, err
}
