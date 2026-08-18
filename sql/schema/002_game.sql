-- +goose Up
CREATE TABLE game (
  date DATE PRIMARY KEY,
  actor_id INT NOT NULL REFERENCES actors(id)
);

-- +goose Down
DROP TABLE game;