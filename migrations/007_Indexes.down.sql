DELETE INDEX users_name_idx ON users;
DELETE INDEX users_id_idx ON users;
DELETE INDEX stats_id_idx ON stats;

DELETE INDEX beatmapsets_id_idx ON beatmapsets;
DELETE INDEX beatmaps_filename_idx ON beatmaps;
DELETE INDEX beatmaps_md5_idx ON beatmaps;
DELETE INDEX beatmaps_id_idx ON beatmaps;

DELETE INDEX idx_score_user_mode_status_pp ON scores;
DELETE INDEX idx_beatmap_mode_status ON scores;
DELETE INDEX idx_beatmap_status ON scores;