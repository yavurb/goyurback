CREATE TYPE post_status AS ENUM ('draft', 'published', 'archived');

CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  title VARCHAR(128) NOT NULL,
  author VARCHAR(64) NOT NULL,
  content TEXT NOT NULL,
  description VARCHAR(255) NOT NULL,
  slug VARCHAR(128) NOT NULL,
  status post_status NOT NULL DEFAULT 'draft',
  published_at TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE projects (
  id SERIAL PRIMARY KEY,
  public_id VARCHAR(15) NOT NULL UNIQUE,
  name VARCHAR(32) NOT NULL,
  description VARCHAR(255) NOT NULL,
  tags VARCHAR[] NOT NULL,
  thumbnail_url VARCHAR(128) NOT NULL,
  website_url VARCHAR(128) NOT NULL,
  live BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  post_id INTEGER,
  FOREIGN KEY (post_id) REFERENCES posts (id)
);
