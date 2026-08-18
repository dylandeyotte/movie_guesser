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

-- name: SelectActor :one
SELECT * FROM actors
ORDER BY RANDOM()
LIMIT 1;

-- name: MarkActor :exec
UPDATE actors 
SET used = TRUE
WHERE name = $1;