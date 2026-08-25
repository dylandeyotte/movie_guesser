-- +goose Up
CREATE TABLE guesses (
  id UUID PRIMARY KEY,
  film_number INT DEFAULT 0 NOT NULL,
  created_at TIMESTAMP NOT NULL,
  date TEXT NOT NULL,
  player_id UUID NOT NULL,
  guess TEXT NOT NUll,
  verdict BOOL NOT NULL
);

-- +goose Down
DROP TABLE guesses;