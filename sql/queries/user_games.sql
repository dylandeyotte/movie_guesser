-- name: CreateUserGame :one
INSERT INTO user_games(id, created_at, date, player_id, actor, correct_guesses, incorrect_guesses, victory)
VALUES (
  gen_random_UUID(),
  NOW(),
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: FetchStats :one
SELECT COUNT(*) AS games_played,
COUNT(*) FILTER (WHERE victory = true AND incorrect_guesses = 0) AS perfect_wins,
COUNT(*) FILTER (WHERE victory = true) AS victories,
COUNT(*) FILTER (WHERE victory = false AND correct_guesses = 0) AS zero_correct,
COUNT(*) FILTER (WHERE victory = false AND correct_guesses = 1) AS one_correct,
COUNT(*) FILTER (WHERE victory = false AND correct_guesses = 2) AS two_correct
FROM user_games WHERE player_id = $1;