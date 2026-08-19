-- +goose Up
CREATE TABLE game (
  date TEXT PRIMARY KEY,
  actor_id INT NOT NULL,
  actor_name TEXT NOT NULL,
  film_1 TEXT NOT NULL,
  film_2 TEXT NOT NULL,
  film_3 TEXT NOT NULL,
  FOREIGN KEY (actor_id)
  REFERENCES actors(id)
);

-- +goose Down
DROP TABLE game;