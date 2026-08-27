-- name: CreateFilm :exec
INSERT INTO films(id, title, poster_path)
VALUES (
  $1,
  $2,
  $3
)
ON CONFLICT (id) DO NOTHING;

-- name: FetchPosters :one
SELECT poster_path FROM films
WHERE id = $1;