
CREATE TABLE beatmapsets
(
    id serial NOT NULL PRIMARY KEY,
    title   character varying(128),
    artist  character varying(128),
    creator character varying(128),
    source  character varying(128),
    tags    character varying(1024) DEFAULT '',
    submission_status int NOT NULL DEFAULT 3,
    has_video boolean NOT NULL DEFAULT false,
    server smallint NOT NULL DEFAULT 0, -- 1: osu! 2: private
    available boolean NOT NULL DEFAULT true,
    submission_date timestamp without time zone NOT NULL DEFAULT now(),
    approved_date timestamp without time zone,
    last_updated timestamp without time zone NOT NULL DEFAULT now(),
    added_at timestamp without time zone DEFAULT now() -- only if server is "osu!"
);

CREATE TABLE beatmaps
(
    id serial NOT NULL PRIMARY KEY,
    set_id int NOT NULL REFERENCES beatmapsets (id),
    mode smallint NOT NULL DEFAULT 0,
    md5 character(32) NOT NULL,
    status smallint NOT NULL DEFAULT 2,
    version  character varying(128) NOT NULL,
    filename character varying(512) NOT NULL,
    submission_date timestamp without time zone NOT NULL DEFAULT now(),
    last_updated timestamp without time zone NOT NULL DEFAULT now(),
    playcount bigint NOT NULL DEFAULT 0,
    passcount bigint NOT NULL DEFAULT 0,
    total_length int NOT NULL,
    max_combo int NOT NULL,
    bpm real NOT NULL DEFAULT 0.00,
    cs real NOT NULL DEFAULT 0.00,
    ar real NOT NULL DEFAULT 0.00,
    od real NOT NULL DEFAULT 0.00,
    hp real NOT NULL DEFAULT 0.00,
    diff real NOT NULL DEFAULT 0.000
);

CREATE TABLE channels
(
    name character varying(32) NOT NULL PRIMARY KEY,
    topic character varying(128) NOT NULL,
    read_permissions int NOT NULL DEFAULT 1,
    write_permissions int NOT NULL DEFAULT 1
);

CREATE TABLE logs
(
    id bigserial NOT NULL PRIMARY KEY,
    level character varying(12) NOT NULL,
    type character varying(250) NOT NULL,
    message character varying NOT NULL,
    "time" timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE messages
(
    id bigserial NOT NULL PRIMARY KEY,
    sender character varying(32) NOT NULL,
    target character varying(32) NOT NULL,
    message character varying(128) NOT NULL,
    "time" time without time zone NOT NULL
);

CREATE TABLE users
(
    id serial NOT NULL PRIMARY KEY,
    name character varying(32) NOT NULL,
    safe_name character varying(32) NOT NULL,
    email character varying(255) NOT NULL,
    pw character(60) NOT NULL, -- bcrypt
    permissions int NOT NULL DEFAULT 1,
    country character varying NOT NULL DEFAULT 'Unknown',
    silence_end timestamp without time zone,
    supporter_end timestamp without time zone,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    latest_activity timestamp without time zone NOT NULL DEFAULT now(),
    restricted boolean NOT NULL DEFAULT false,
    activated boolean NOT NULL DEFAULT false,
    preferred_mode int NOT NULL DEFAULT 0,
    playstyle int NOT NULL DEFAULT 0,
    UNIQUE(name, safe_name, email)
);

CREATE TABLE stats
(
    id int NOT NULL REFERENCES users (id),
    mode smallint NOT NULL, -- osu!, taiko, ctb and mania
    rank int NOT NULL DEFAULT 0,
    tscore bigint NOT NULL DEFAULT 0,
    rscore bigint NOT NULL DEFAULT 0,
    pp real NOT NULL DEFAULT 0,
    playcount bigint NOT NULL DEFAULT 0,
    playtime int NOT NULL DEFAULT 0,
    acc real NOT NULL DEFAULT 0.000,
    max_combo int NOT NULL DEFAULT 0,
    total_hits bigint NOT NULL DEFAULT 0,
    replay_views bigint NOT NULL DEFAULT 0,
    xh_count int NOT NULL DEFAULT 0,
    x_count int NOT NULL DEFAULT 0,
    sh_count int NOT NULL DEFAULT 0,
    s_count int NOT NULL DEFAULT 0,
    a_count int NOT NULL DEFAULT 0,
    b_count int NOT NULL DEFAULT 0,
    c_count int NOT NULL DEFAULT 0,
    d_count int NOT NULL DEFAULT 0,
    CONSTRAINT stats_pkey PRIMARY KEY (id, mode)
);

CREATE TABLE screenshots
(
    id bigserial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users (id),
    created_at time without time zone NOT NULL DEFAULT now(),
    hidden boolean NOT NULL DEFAULT false,
    content bytea NOT NULL
);

CREATE TABLE scores
(
    id bigserial NOT NULL PRIMARY KEY,
    beatmap_id int NOT NULL REFERENCES beatmaps (id),
    user_id int NOT NULL REFERENCES users (id),
    client_version int NOT NULL,
    client_hash character varying NOT NULL,
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
    replay bytea,
    screenshot bytea,
    processes character varying,
    failtime int
);

-- Status:
-- 0: Friend
-- 1: Blocked
CREATE TABLE relationships
(
    user_id int NOT NULL REFERENCES users (id),
    target_id int NOT NULL REFERENCES users (id),
    status smallint NOT NULL,
    PRIMARY KEY (user_id, target_id)
);

CREATE TABLE ratings
(
    user_id int NOT NULL REFERENCES users (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    map_checksum character(32) NOT NULL,
    rating smallint NOT NULL,
    PRIMARY KEY (user_id, map_checksum)
);

CREATE TABLE plays
(
    user_id int NOT NULL REFERENCES users (id),
    beatmap_id int NOT NULL REFERENCES beatmaps (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    count int NOT NULL,
    beatmap_file character varying NOT NULL,
    PRIMARY KEY (user_id, beatmap_id)
);

create table favourites
(
	user_id int NOT NULL REFERENCES users (id),
	set_id int NOT NULL REFERENCES beatmapsets (id),
	created_at time without time zone NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, set_id)
);

CREATE TABLE achievements
(
    user_id int NOT NULL REFERENCES users (id),
    name character varying NOT NULL,
    category character varying NOT NULL,
    filename character varying NOT NULL,
    unlocked_at time without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, name)
);

CREATE TABLE comments
(
    id bigserial NOT NULL,
    target_id int NOT NULL,
    target_type character varying(6) NOT NULL,
    user_id int NOT NULL REFERENCES users (id),
    "time" int NOT NULL,
    comment character varying(80) NOT NULL,
    format character varying(10),
    mode smallint NOT NULL DEFAULT 0,
    CONSTRAINT comments_pkey PRIMARY KEY (id)
);

INSERT INTO users (name, safe_name, email, pw, permissions, country, activated)
VALUES ('BanchoBot', 'banchobot', 'bot@example.com', '------------------------------------------------------------', 21, 'Oceania', true),
       ('peppy', 'peppy', 'pe@ppy.sh', '$2b$12$W5ppLwlSEJ3rpJQRq8UcX.QA5cTm7HvsVpn6MXQHE/6OEO.Iv4DGW', 21, 'Australia', true);

INSERT INTO stats (id, mode)
VALUES (1, 0),
       (1, 1),
       (1, 2),
       (1, 3);

INSERT INTO channels (name, topic, read_permissions, write_permissions)
VALUES ('#osu', 'General discussion.', 1, 1),
       ('#highlight', 'Public announcements.', 1, 8),
       ('#lobby', 'Multiplayer lobby discussion room.', 1, 1),
       ('#admin', 'General discussion for administrators.', 16, 16);
