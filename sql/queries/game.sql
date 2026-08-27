-- name: CreateGame :one
INSERT INTO game(date, created_at, actor_id, actor_name, film_1, film_2, film_3, film_1_id, film_2_id, film_3_id)
VALUES (
  $1,
  NOW(),
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9
)
RETURNING *;

-- name: ReturnGame :one
SELECT * FROM game
WHERE date = $1;