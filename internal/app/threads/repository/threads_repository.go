package repository

import (
	"BD-v2/internal/app/threads/models"
	"context"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type ThreadsRepository struct {
	DBPool *pgxpool.Pool
}

func NewThreadsRepository(db *pgxpool.Pool) *ThreadsRepository {
	return &ThreadsRepository{
		DBPool: db,
	}
}

func (rep *ThreadsRepository) CreateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	var query string
	var err error
	if thread.Slug != "" {
		if thread.Created != "" {
			query = `insert into threads (title, author, forum, message, slug, created)
			values ($1, $2, $3, $4, $5, $6) returning id, title, author, forum,
			message, votes, slug`
			err = rep.DBPool.QueryRow(ctx, query, thread.Title, thread.Author,
				thread.Forum, thread.Message, thread.Slug, thread.Created).Scan(&thread.ID, &thread.Title,
				&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
				&thread.Slug,
			)
		} else {
			query = `insert into threads (title, author, forum, message, slug)
			values ($1, $2, $3, $4, $5) returning id, title, author, forum,
			message, votes, slug`
			err = rep.DBPool.QueryRow(ctx, query, thread.Title, thread.Author,
				thread.Forum, thread.Message, thread.Slug).Scan(&thread.ID, &thread.Title,
				&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
				&thread.Slug)
		}
	} else {
		if thread.Created != "" {
			query = `insert into threads (title, author, forum, message, created)
			values ($1, $2, $3, $4, $5) returning id, title, author, forum,
			message, votes`
			err = rep.DBPool.QueryRow(ctx, query, thread.Title, thread.Author,
				thread.Forum, thread.Message, &thread.Created).Scan(&thread.ID, &thread.Title,
				&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
			)
		} else {
			query = `insert into threads (title, author, forum, message)
			values ($1, $2, $3, $4) returning id, title, author, forum,
			message, votes`
			err = rep.DBPool.QueryRow(ctx, query, thread.Title, thread.Author,
				thread.Forum, thread.Message).Scan(&thread.ID, &thread.Title,
				&thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
			)
		}
	}
	return thread, err
}

func (rep *ThreadsRepository) FindThreadSlug(ctx context.Context, slug string) (*models.Thread, error) {
	query := `select id, title, author, forum,
	message, votes, slug, created from threads
	where slug = $1`
	thread := &models.Thread{}
	date := time.Time{}

	err := rep.DBPool.QueryRow(ctx, query, slug).Scan(&thread.ID, &thread.Title, &thread.Author,
		&thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &date,
	)
	thread.Created = strfmt.DateTime(date.UTC()).String()
	return thread, err
}
