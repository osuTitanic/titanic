-- This column now stores the entire client version string, instead of just the numeric version
ALTER TABLE scores ADD COLUMN client_version_string varchar(25) NOT NULL DEFAULT '';

-- Create index to speed up the backfill query
CREATE INDEX IF NOT EXISTS idx_logins_user_time ON logins (user_id, "time" DESC);
ANALYZE logins;

-- Backfill existing score metadata
-- We can do this by retrieving the osu_version from
-- the most recent login from when the score was set
UPDATE scores s
SET client_version_string = COALESCE(
    (
        SELECT l.osu_version FROM logins l
        WHERE l.user_id = s.user_id AND l."time" <= s.submitted_at AND l.osu_version != 'irc'
        ORDER BY l."time" DESC LIMIT 1
    ), ''
);