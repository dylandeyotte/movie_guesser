-- name: InsertActor :one
INSERT INTO actors(id, name, used)
VALUES (
  $1,
  $2,
  FALSE
)
RETURNING *;

-- name: RemoveActor :exec
DELETE FROM actors
WHERE name = $1;

-- name: ReturnActors :many
SELECT * FROM actors;