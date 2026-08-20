-- +goose Up
CREATE TABLE guesses (
  id UUID PRIMARY KEY,
  date TEXT NOT NULL,
  player_id UUID NOT NULL,
  guess TEXT NOT NUll,
  verdict BOOL NOT NULL
);

-- +goose Down
DROP TABLE guesses;