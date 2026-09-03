-- +goose Up
CREATE TABLE user_games (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  date TEXT NOT NULL,
  player_id UUID NOT NULL,
  actor TEXT NOT NULL,
  correct_guesses INT NOT NULL,
  incorrect_guesses INT NOT NULL,
  victory BOOL NOT NULL
);

-- +goose Down
DROP TABLE user_games;
