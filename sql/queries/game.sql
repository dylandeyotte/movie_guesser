-- name: CreateGame :one
INSERT INTO game(date, actor_id, actor_name, film_1, film_2, film_3)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: ReturnGame :one
SELECT * FROM game
WHERE date = $1;