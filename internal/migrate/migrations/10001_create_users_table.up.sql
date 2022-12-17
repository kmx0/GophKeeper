CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  login    varchar(255) UNIQUE,
  password   varchar(255),
  created_at timestamp
)