-- +goose Up
SELECT 'up SQL query';
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    category_name TEXT NOT NULL,
    info TEXT
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE categories;