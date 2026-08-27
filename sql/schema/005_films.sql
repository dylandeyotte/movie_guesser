-- +goose Up 
CREATE TABLE films (
  id INT PRIMARY KEY,
  title TEXT NOT NULL,
  poster_path TEXT NOT NULL
);

-- +goose Down
DROP TABLE films;