-- +goose Up
ALTER TABLE game
ADD COLUMN film_1_id INT NOT NULL, 
ADD COLUMN film_2_id INT NOT NULL, 
ADD COLUMN film_3_id INT NOT NULL;

-- +goose Down
ALTER TABLE game
DROP COLUMN film_1_id, film_2_id, film_3_id;