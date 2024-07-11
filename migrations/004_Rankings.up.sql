CREATE TABLE IF NOT EXISTS scores
(
    id bigserial NOT NULL PRIMARY KEY,
    beatmap_id int NOT NULL REFERENCES beatmaps (id),
    user_id int NOT NULL REFERENCES users (id),
    client_version int NOT NULL,
    client_hash character varying,
    score_checksum character(32) NOT NULL,
    mode smallint NOT NULL,
    pp real NOT NULL,
    acc real NOT NULL,
    total_score bigint NOT NULL,
    max_combo int NOT NULL,
    mods int NOT NULL,
    perfect boolean NOT NULL,
    n300 int NOT NULL,
    n100 int NOT NULL,
    n50 int NOT NULL,
    nmiss int NOT NULL,
    ngeki int NOT NULL,
    nkatu int NOT NULL,
    grade character varying(2) NOT NULL DEFAULT 'N',
    status smallint NOT NULL,
    submitted_at timestamp with time zone NOT NULL DEFAULT now(),
	replay_md5 CHARACTER(32),
    processes character varying,
    failtime int,
    pinned boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS achievements
(
    user_id int NOT NULL REFERENCES users (id),
    name character varying NOT NULL,
    category character varying NOT NULL,
    filename character varying NOT NULL,
    unlocked_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, name)
);

CREATE TABLE IF NOT EXISTS benchmarks
(
    id serial NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    smoothness real NOT NULL,
    framerate int NOT NULL,
    score int NOT NULL,
    grade character varying(2) NOT NULL DEFAULT 'N',
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_score_user_mode_status_pp ON scores (user_id, mode, status, pp DESC);
CREATE INDEX IF NOT EXISTS idx_beatmap_mode_status ON scores (beatmap_id, mode, status);
CREATE INDEX IF NOT EXISTS idx_achievement_user_id ON achievements (user_id);
CREATE INDEX IF NOT EXISTS idx_beatmap_status ON scores (status);