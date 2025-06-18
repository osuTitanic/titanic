BEGIN;

-- Add new columns
ALTER TABLE profile_activity
    ADD COLUMN data jsonb NOT NULL DEFAULT '{}',
    ADD COLUMN type smallint NOT NULL DEFAULT 0,
    ADD COLUMN hidden boolean NOT NULL DEFAULT false;

-- Migrate existing rows partially
UPDATE profile_activity
SET
    type = CASE
        WHEN activity_text LIKE '%has risen%'                 THEN 1 -- ranks_gained
        WHEN activity_text LIKE '%taken the lead%'            THEN 2 -- number_one
        WHEN activity_text LIKE '%achieved rank#%'            THEN 3 -- beatmap_leaderboard_rank
        WHEN activity_text LIKE '%has lost first place%'      THEN 4 -- lost_first_place
        WHEN activity_text LIKE '%has set the new pp record%' THEN 5 -- pp_record
        WHEN activity_text LIKE '%got a new top play%'        THEN 6 -- top_play
        WHEN activity_text LIKE '%unlocked an achievement%'   THEN 7 -- achievement_unlocked
    ELSE type
END;

-- If no profile_activity rows exist, delete the old columns
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM profile_activity) THEN
        EXECUTE 'ALTER TABLE profile_activity
                 DROP COLUMN activity_text,
                 DROP COLUMN activity_args,
                 DROP COLUMN activity_links';
    END IF;
END
$$;

COMMIT;