package postgres

import "context"

const createLonger = `--name: CreateLonger :one
INSERT INTO urls (id, original)
VALUES ($1, $2) RETURNING *
`

// CreateLongerParams struct
type CreateLongerParams struct {
	ID       string `json:"id"`
	Original string `json:"original"`
}

func (q *Queries) CreateLonger(ctx context.Context, arg CreateLongerParams) (URL, error) {
	row := q.db.QueryRowContext(ctx, createLonger, arg.ID, arg.Original)
	var i URL
	err := row.Scan(
		&i.ID,
		&i.Original,
	)
	return i, err
}

const getOriginal = `-- name: GetOriginal :one
SELECT original FROM urls
WHERE id = $1
`

func (q *Queries) GetOriginal(ctx context.Context, id string) (URL, error) {
	row := q.db.QueryRowContext(ctx, getOriginal, id)
	var i URL
	i.ID = id
	err := row.Scan(
		&i.Original,
	)
	return i, err
}

const getUrls = `--name: GetAll :many
SELECT id, original FROM urls
`

func (q *Queries) GetUrls(ctx context.Context) ([]URL, error) {
	rows, err := q.db.QueryContext(ctx, getUrls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []URL
	for rows.Next() {
		var u URL
		if err := rows.Scan(&u.ID, &u.Original); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, err
}
