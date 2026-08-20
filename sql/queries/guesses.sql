-- name: CreateGuess :one
INSERT INTO guesses(id, date, player_id, guess, verdict)
VALUES (
  gen_random_UUID(),
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: StrikeCount :one
SELECT COUNT(*) FROM guesses
WHERE date = $1
AND player_id = $2
AND verdict = FALSE;