DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS posts;

DROP TABLE IF EXISTS comments;

DROP INDEX IF EXISTS idx_comments_post_id;

DROP INDEX IF EXISTS idx_comments_parent_comment_id;