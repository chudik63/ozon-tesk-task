CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  username VARCHAR(10) NOT NULL UNIQUE,
  email VARCHAR(80) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS posts (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  user_id BIGINT REFERENCES users(id),
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  comments_allowed BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  post_id BIGINT REFERENCES posts(id),
  user_id BIGINT REFERENCES users(id),
  parent_comment_id BIGINT REFERENCES comments(id),
  content TEXT NOT NULL CHECK (char_length(content) <= 2000),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments USING btree (post_id);

CREATE INDEX IF NOT EXISTS idx_comments_parent_comment_id ON comments USING btree (parent_comment_id);