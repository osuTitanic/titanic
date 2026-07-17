ALTER TABLE forum_posts
ADD COLUMN smilies_disabled boolean NOT NULL DEFAULT false;

-- Old forum posts should have smilies disabled by default
UPDATE forum_posts SET smilies_disabled = true
