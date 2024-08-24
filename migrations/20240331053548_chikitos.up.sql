CREATE TABLE chikitos (
  id SERIAL PRIMARY KEY,
  public_id VARCHAR(15) NOT NULL UNIQUE,
  url VARCHAR(255) NOT NULL,
  description VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);