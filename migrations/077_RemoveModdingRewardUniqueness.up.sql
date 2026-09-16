BEGIN;

ALTER TABLE beatmap_modding
    DROP CONSTRAINT IF EXISTS uq_modding_reward;

COMMIT;
