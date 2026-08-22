-- name: CreateGuess :one
INSERT INTO guesses(id, created_at, date, player_id, guess, verdict)
VALUES (
  gen_random_UUID(),
  NOW(),
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: FetchGuess :one
SELECT * FROM guesses
WHERE date = $1
AND player_id = $2
AND guess = $3;

-- name: FetchGuessList :many
SELECT date, guess, verdict FROM guesses
WHERE date = $1
AND player_id = $2;

-- name: StrikeCount :one
SELECT COUNT(*) FROM guesses
WHERE date = $1
AND player_id = $2
AND verdict = FALSE;