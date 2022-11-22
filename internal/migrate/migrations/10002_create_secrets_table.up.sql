CREATE TABLE secrets (
  id SERIAL PRIMARY KEY,
  users_id integer,
  key    VARCHAR(255),
  type varchar(255),
  value   TEXT,
  created_at TIMESTAMP,
  CONSTRAINT fk_users
	FOREIGN KEY(users_id)
	REFERENCES users(id)
	ON DELETE CASCADE
);
