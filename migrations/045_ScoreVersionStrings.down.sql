ALTER TABLE scores DROP COLUMN client_version_string;
DROP INDEX IF EXISTS idx_logins_user_time;