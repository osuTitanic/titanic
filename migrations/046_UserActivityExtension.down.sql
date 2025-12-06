DROP VIEW user_count;
ALTER TABLE user_activity RENAME TO user_count;
ALTER TABLE user_count RENAME COLUMN osu_count TO count;
ALTER TABLE user_count DROP COLUMN irc_count;
ALTER TABLE user_count DROP COLUMN mp_count;