-- Rename table, but keep a backup view for compatibility
ALTER TABLE user_count RENAME TO user_activity;
ALTER TABLE user_activity RENAME COLUMN count TO osu_count;
ALTER TABLE user_activity ADD COLUMN irc_count INTEGER DEFAULT 0;
ALTER TABLE user_activity ADD COLUMN mp_count INTEGER DEFAULT 0;

CREATE VIEW user_count AS SELECT time, osu_count AS count FROM user_activity;