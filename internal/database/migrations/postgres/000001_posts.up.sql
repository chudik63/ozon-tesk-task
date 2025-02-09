-- CREATE TABLE IF NOT EXISTS users (
--   id SERIAL PRIMARY KEY,
--   username VARCHAR(10) NOT NULL UNIQUE,
--   email VARCHAR(80) NOT NULL UNIQUE
-- );

CREATE TABLE IF NOT EXISTS posts (
  id SERIAL PRIMARY KEY,
  user_id INTEGER,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  comments_allowed BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
  -- FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS comments (
  id SERIAL PRIMARY KEY,
  post_id INTEGER,
  user_id INTEGER,
  parent_comment_id INTEGER REFERENCES comments(id),
  content VARCHAR(2000) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id),
  -- FOREIGN KEY (user_id) REFERENCES users(id),
);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_comment_id ON comments (parent_comment_id);