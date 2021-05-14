package repository

import (
	models2 "BD-v2/internal/app/posts/models"
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

func (rep *ThreadsRepository) FindThreadID(ctx context.Context, threadID int) (*models.Thread, error) {
	query := `select id, title, author, forum,
	message, votes, slug, created from threads
	where id = $1`

	thread := &models.Thread{}
	date := time.Time{}

	err := rep.DBPool.QueryRow(ctx, query, threadID).Scan(&thread.ID, &thread.Title, &thread.Author,
		&thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &date,
	)
	thread.Created = strfmt.DateTime(date.UTC()).String()
	return thread, err
}

func (rep *ThreadsRepository) CreatePost(ctx context.Context, post *models2.Post) (*models2.Post, error) {
	query := `insert into posts (parent, author, message, forum, thread)
		values (NULLIF($1, 0), $2, $3, $4, $5) returning id, parent, author, message, 
		is_edited, forum, thread, created`
	date := time.Time{}

	err := rep.DBPool.QueryRow(ctx, query, post.Parent, post.Author, post.Message,
		post.Forum, post.Thread).Scan(&post.ID, &post.Parent, &post.Author, &post.Message,
		&post.ISEdited, &post.Forum, &post.Thread, &date)

	post.Created = strfmt.DateTime(date.UTC()).String()
	return post, err
}

func (rep *ThreadsRepository) AddVoice(ctx context.Context, voice *models.Vote) error {
	query := `insert into thread_votes (nickname, voice, thread_id) VALUES ($1, $2, $3)`

	_, err := rep.DBPool.Exec(ctx, query, voice.Nickname, voice.Voice, voice.ThreadID)
	return err
}

func (rep *ThreadsRepository) GetVoice(ctx context.Context, voice *models.Vote) (*models.Vote, error) {
	query := `select nickname, voice, thread_id from thread_votes
		where nickname = $1 and thread_id = $2`

	err := rep.DBPool.QueryRow(ctx, query, voice.Nickname, voice.ThreadID).Scan(&voice.Nickname,
		&voice.Voice, &voice.ThreadID)
	return voice, err
}

func (rep *ThreadsRepository) UpdateVoice(ctx context.Context, voice *models.Vote) error {
	query := `update thread_votes
		set voice = $1
		where nickname = $2 and thread_id = $3`

	_, err := rep.DBPool.Exec(ctx, query, voice.Voice, voice.Nickname, voice.ThreadID)
	return err
}
