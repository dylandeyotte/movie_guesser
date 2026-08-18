-- name: CreateGame :one
INSERT INTO game(date, actor_id)
VALUES (
  $1,
  $2
)
RETURNING *;