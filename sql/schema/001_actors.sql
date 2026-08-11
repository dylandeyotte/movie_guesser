-- +goose Up
CREATE TABLE actors (
  id INT PRIMARY KEY,
  name TEXT NOT NULL,
  used BOOL NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE actors;