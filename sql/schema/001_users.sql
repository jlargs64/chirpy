-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL,
  createdAt TIMESTAMP NOT NULL,
  updatedAt TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE users;
