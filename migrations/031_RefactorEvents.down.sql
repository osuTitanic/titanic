ALTER TABLE profile_activity
    DROP COLUMN IF EXISTS data;
ALTER TABLE profile_activity
    DROP COLUMN IF EXISTS type;

ALTER TABLE profile_activity
    ADD COLUMN IF NOT EXISTS activity_text TEXT;
ALTER TABLE profile_activity
    ADD COLUMN IF NOT EXISTS activity_args TEXT;
ALTER TABLE profile_activity
    ADD COLUMN IF NOT EXISTS activity_links TEXT;