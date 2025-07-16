ALTER TABLE users ADD COLUMN avatar_hash CHARACTER(32) DEFAULT NULL;
ALTER TABLE users ADD COLUMN avatar_last_changed timestamp without time zone NOT NULL DEFAULT now();