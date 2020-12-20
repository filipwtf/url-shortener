--name: GetAll :many
SELECT id, original FROM urls;

-- name: GetOriginal :one
SELECT original FROM urls
WHERE id = $1;

--name: CreateLonger :one
INSERT INTO urls (id, original)
VALUES ($1, $2) RETURNING *;
