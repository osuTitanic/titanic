CREATE TABLE IF NOT EXISTS beatmap_collaboration
(
    user_id int NOT NULL REFERENCES users (id),
    beatmap_id int NOT NULL REFERENCES beatmaps (id),
    is_beatmap_author boolean NOT NULL DEFAULT false,
    allow_resource_updates boolean NOT NULL DEFAULT false,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, beatmap_id)
);

CREATE TABLE IF NOT EXISTS beatmap_collaboration_requests
(
    id serial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users (id),
    target_id int NOT NULL REFERENCES users (id),
    beatmap_id int NOT NULL REFERENCES beatmaps (id),
    allow_resource_updates boolean NOT NULL DEFAULT false,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS beatmap_collaboration_blacklist
(
    user_id int NOT NULL REFERENCES users (id),
    target_id int NOT NULL REFERENCES users (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, target_id)
);

CREATE INDEX IF NOT EXISTS idx_beatmap_collaboration_user_id ON beatmap_collaboration (user_id);
CREATE INDEX IF NOT EXISTS idx_beatmap_collaboration_beatmap_id ON beatmap_collaboration (beatmap_id);
CREATE INDEX IF NOT EXISTS idx_beatmap_collaboration_requests_user_id ON beatmap_collaboration_requests (user_id);
CREATE INDEX IF NOT EXISTS idx_beatmap_collaboration_requests_target_id ON beatmap_collaboration_requests (target_id);
CREATE INDEX IF NOT EXISTS idx_beatmap_collaboration_requests_beatmap_id ON beatmap_collaboration_requests (beatmap_id);
CREATE INDEX IF NOT EXISTS idx_beatmap_collaboration_blacklist_user_id ON beatmap_collaboration_blacklist (user_id, target_id);