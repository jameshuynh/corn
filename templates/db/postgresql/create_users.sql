-- +mig Up
CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  encrypted_password VARCHAR(255) DEFAULT '',
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +mig Down
DROP TABLE users;
