BEGIN;

ALTER TABLE beatmap_modding
    ADD CONSTRAINT uq_modding_reward
    UNIQUE (target_id, sender_id, set_id, post_id);

COMMIT;
