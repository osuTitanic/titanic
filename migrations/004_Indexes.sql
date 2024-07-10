CREATE INDEX users_name_idx ON users (name);
CREATE INDEX users_id_idx ON users (id);
CREATE INDEX stats_id_idx ON stats (id);

CREATE INDEX beatmapsets_id_idx ON beatmapsets (id);
CREATE INDEX beatmaps_filename_idx ON beatmaps (filename);
CREATE INDEX beatmaps_md5_idx ON beatmaps (md5);
CREATE INDEX beatmaps_id_idx ON beatmaps (id);

CREATE INDEX idx_score_user_mode_status_pp ON scores (user_id, mode, status, pp DESC);
CREATE INDEX idx_beatmap_mode_status ON scores (beatmap_id, mode, status);
CREATE INDEX idx_beatmap_status ON scores (status);