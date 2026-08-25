-- name: CreateGuess :one
INSERT INTO guesses(id, film_number, created_at, date, player_id, guess, verdict)
VALUES (
  gen_random_UUID(),
  $1,
  NOW(),
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: FetchGuess :one
SELECT * FROM guesses
WHERE date = $1
AND player_id = $2
AND guess = $3;

-- name: FetchGuessList :many
SELECT date, film_number, guess, verdict FROM guesses
WHERE date = $1
AND player_id = $2;

-- name: StrikeCount :one
SELECT COUNT(*) FROM guesses
WHERE date = $1
AND player_id = $2
AND verdict = FALSE;