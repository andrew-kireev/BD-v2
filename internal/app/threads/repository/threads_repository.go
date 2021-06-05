package repository

import (
	models2 "BD-v2/internal/app/posts/models"
	"BD-v2/internal/app/threads/models"
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
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

func (rep *ThreadsRepository) CreatePost(posts []*models2.Post, thread *models.Thread) ([]*models2.Post, error) {
	query := `INSERT INTO posts(author, created, message, parent, thread, forum) VALUES `
	empty := make([]*models2.Post, 0)
	if len(posts) == 0 {
		return empty, nil
	}

	timeCreated := time.Now()
	valuesNames := make([]string, 0)
	var values []interface{}
	i := 1
	for _, element := range posts {
		valuesNames = append(valuesNames, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d)",
			i, i+1, i+2, i+3, i+4, i+5))
		i += 6
		values = append(values, element.Author, timeCreated, element.Message, element.Parent, thread.ID, thread.Forum)
	}

	query += strings.Join(valuesNames, ",")
	query += " RETURNING author, created, forum, id, is_edited, message, parent, thread"
	row, err := rep.DBPool.Query(context.Background(), query, values...)

	if err != nil {
		return empty, err
	}
	defer func() {
		if row != nil {
			row.Close()
		}
	}()

	for row.Next() {
		post := &models2.Post{}
		var created time.Time
		err = row.Scan(&post.Author, &created, &post.Forum, &post.ID, &post.ISEdited,
			&post.Message, &post.Parent, &post.Thread)

		if err != nil {
			return empty, err
		}
		post.Created = strfmt.DateTime(created.UTC()).String()
		empty = append(empty, post)
	}
	return empty, err
}

func (rep *ThreadsRepository) GetPosts(threadID, limit, since int, desc bool) ([]*models2.Post, error) {
	query := `SELECT id, author, message, is_edited, forum, thread, created, parent FROM posts WHERE 
	thread = $1 `

	if desc {
		if since > 0 {
			query += fmt.Sprintf("AND id < %d ", since)
		}
		query += `ORDER BY id DESC `
	} else {
		if since > 0 {
			query += fmt.Sprintf("AND id > %d ", since)
		}
		query += `ORDER BY id `
	}
	query += `LIMIT NULLIF($2, 0)`
	posts := make([]*models2.Post, 0)

	row, err := rep.DBPool.Query(context.Background(), query, threadID, limit)

	if err != nil {
		return posts, err
	}
	defer func() {
		if row != nil {
			row.Close()
		}
	}()

	for row.Next() {
		post := &models2.Post{}
		created := &time.Time{}

		err = row.Scan(&post.ID, &post.Author, &post.Message, &post.ISEdited, &post.Forum,
			&post.Thread, created, &post.Parent)

		if err != nil {
			return posts, err
		}
		post.Created = strfmt.DateTime(created.UTC()).String()
		posts = append(posts, post)
	}
	return posts, err
}

func (rep *ThreadsRepository) GetPostsTree(threadID, limit, since int, desc bool) ([]*models2.Post, error) {
	var query string
	sinceQuery := ""
	if since != 0 {
		if desc {
			sinceQuery = `AND path < `
		} else {
			sinceQuery = `AND path > `
		}
		sinceQuery += fmt.Sprintf(`(SELECT path FROM posts WHERE id = %d)`, since)
	}
	if desc {
		query = fmt.Sprintf(
			`SELECT id, author, message, is_edited, forum, thread, created, parent FROM posts WHERE thread=$1 %s ORDER BY path DESC, id DESC LIMIT NULLIF($2, 0);`, sinceQuery)
	} else {
		query = fmt.Sprintf(
			`SELECT id, author, message, is_edited, forum, thread, created, parent FROM posts WHERE thread=$1 %s ORDER BY path, id LIMIT NULLIF($2, 0);`, sinceQuery)
	}
	var posts []*models2.Post
	row, err := rep.DBPool.Query(context.Background(),query, threadID, limit)

	if err != nil {
		return posts, err
	}
	defer func() {
		if row != nil {
			row.Close()
		}
	}()

	for row.Next() {
		post := &models2.Post{}
		created := &time.Time{}

		err = row.Scan(&post.ID, &post.Author, &post.Message, &post.ISEdited, &post.Forum,
			&post.Thread, created, &post.Parent)

		if err != nil {
			return posts, err
		}
		post.Created = strfmt.DateTime(created.UTC()).String()
		posts = append(posts, post)

	}
	return posts, err
}

func (rep *ThreadsRepository) GetPostsParentTree(threadID, limit, since int, desc bool) ([]*models2.Post, error) {
	var query string
	sinceQuery := ""
	if since != 0 {
		if desc {
			sinceQuery = `AND PATH[1] < `
		} else {
			sinceQuery = `AND PATH[1] > `
		}
		sinceQuery += fmt.Sprintf(`(SELECT path[1] FROM posts WHERE id = %d)`, since)
	}

	parentsQuery := fmt.Sprintf(
		`SELECT id FROM posts WHERE thread = $1 AND parent IS NULL %s`, sinceQuery)

	if desc {
		parentsQuery += `ORDER BY id DESC`
		if limit > 0 {
			parentsQuery += fmt.Sprintf(` LIMIT %d`, limit)
		}
		query = fmt.Sprintf(
			`SELECT id, author, message, is_edited, forum, thread, created, parent FROM posts WHERE path[1] IN (%s) ORDER BY path[1] DESC, path, id;`, parentsQuery)
	} else {
		parentsQuery += `ORDER BY id`
		if limit > 0 {
			parentsQuery += fmt.Sprintf(` LIMIT %d`, limit)
		}
		query = fmt.Sprintf(
			`SELECT id, author, message, is_edited, forum, thread, created, parent FROM posts WHERE path[1] IN (%s) ORDER BY path,id;`, parentsQuery)
	}
	var posts []*models2.Post
	row, err := rep.DBPool.Query(context.Background(), query, threadID)

	if err != nil {
		return posts, err
	}

	defer func() {
		if row != nil {
			row.Close()
		}
	}()

	for row.Next() {
		post := &models2.Post{}
		created := &time.Time{}

		err = row.Scan(&post.ID, &post.Author, &post.Message, &post.ISEdited, &post.Forum,
			&post.Thread, created, &post.Parent)

		if err != nil {
			return posts, err
		}
		post.Created = strfmt.DateTime(created.UTC()).String()
		posts = append(posts, post)

	}
	return posts, err
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

func (rep *ThreadsRepository) UpdateVoice(ctx context.Context, voice *models.Vote) (int, error) {
	query := `update thread_votes
		set voice = $1
		where nickname = $2 and thread_id = $3 and voice != $1`

	res, err := rep.DBPool.Exec(ctx, query, voice.Voice, voice.Nickname, voice.ThreadID)
	if err != nil {
		return 0, err
	}
	return int(res.RowsAffected()), err
}

func (rep *ThreadsRepository) UpdateThreadID(ctx context.Context, thread *models.Thread) error {
	query := `update threads set title = $1, message = $2
	where id = $3
`

	_, err := rep.DBPool.Exec(ctx, query, thread.Title, thread.Message, thread.ID)
	return err
}

func (rep *ThreadsRepository) UpdateThreadSlug(ctx context.Context, thread *models.Thread) error {
	query := `update threads set title = $1, message = $2
	where slug = $3
`
	_, err := rep.DBPool.Exec(ctx, query, thread.Title, thread.Message, thread.Slug)
	return err
}
