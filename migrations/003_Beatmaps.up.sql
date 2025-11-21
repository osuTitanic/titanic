CREATE TABLE IF NOT EXISTS beatmapsets
(
    id serial NOT NULL PRIMARY KEY,
    title character varying(128),
    title_unicode character varying(128),
    artist character varying(128),
    artist_unicode character varying(128),
    source character varying(128),
    source_unicode character varying(128),
    creator character varying(128),
    tags character varying(1024) DEFAULT '',
    display_title text NOT NULL DEFAULT '',
    description text DEFAULT '',
    submission_status int NOT NULL DEFAULT 3,
    has_video boolean NOT NULL DEFAULT false,
    has_storyboard boolean NOT NULL DEFAULT false,
    server smallint NOT NULL DEFAULT 0, -- 0: ppy 1: private
    topic_id int REFERENCES forums (id) DEFAULT NULL, -- only if server is "private"
    creator_id int REFERENCES users (id) DEFAULT NULL, -- only if server is "private"
    available boolean NOT NULL DEFAULT true,
    submission_date timestamp without time zone NOT NULL DEFAULT now(),
    approved_date timestamp without time zone,
    approved_by int REFERENCES users (id) DEFAULT NULL,
    last_updated timestamp without time zone NOT NULL DEFAULT now(),
    added_at timestamp without time zone DEFAULT now(),
    osz_filesize int NOT NULL DEFAULT 0,
    osz_filesize_novideo int NOT NULL DEFAULT 0,
    language_id smallint NOT NULL DEFAULT 1,
    genre_id smallint NOT NULL DEFAULT 1,
    star_priority int NOT NULL DEFAULT 0,
    "offset" int NOT NULL DEFAULT 0,
    meta_hash character(32) DEFAULT NULL,
    info_hash character(32) DEFAULT NULL,
    body_hash character(32) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS beatmaps
(
    id serial NOT NULL PRIMARY KEY,
    set_id int NOT NULL REFERENCES beatmapsets (id),
    mode smallint NOT NULL DEFAULT 0,
    md5 character(32) NOT NULL,
    status smallint NOT NULL DEFAULT 2,
    version character varying(128) NOT NULL,
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

CREATE TABLE IF NOT EXISTS beatmap_nominations
(
    user_id int NOT NULL REFERENCES users (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, set_id)
);

CREATE TABLE IF NOT EXISTS beatmap_modding
(
    id serial NOT NULL PRIMARY KEY,
    target_id int NOT NULL REFERENCES users (id),
    sender_id int NOT NULL REFERENCES users (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    post_id int NOT NULL REFERENCES forum_posts (id),
    amount int NOT NULL DEFAULT 0,
    "time" timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ratings
(
    user_id int NOT NULL REFERENCES users (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    map_checksum character(32) NOT NULL,
    rating smallint NOT NULL,
    PRIMARY KEY (user_id, map_checksum)
);

CREATE TABLE IF NOT EXISTS plays
(
    user_id int NOT NULL REFERENCES users (id),
    beatmap_id int NOT NULL REFERENCES beatmaps (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    count int NOT NULL,
    beatmap_file character varying NOT NULL,
    PRIMARY KEY (user_id, beatmap_id)
);

CREATE TABLE IF NOT EXISTS favourites
(
	user_id int NOT NULL REFERENCES users (id),
	set_id int NOT NULL REFERENCES beatmapsets (id),
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, set_id)
);

CREATE TABLE IF NOT EXISTS comments
(
    id bigserial NOT NULL,
    target_id int NOT NULL,
    target_type character varying(6) NOT NULL,
    user_id int NOT NULL REFERENCES users (id),
    "time" int NOT NULL,
    comment character varying(80) NOT NULL,
    format character varying(10),
    mode smallint NOT NULL DEFAULT 0,
	color character varying(8),
    CONSTRAINT comments_pkey PRIMARY KEY (id)
);

-- resource type:
--   0: osz
--   1: osz_novideo
--   2: beatmap_file
--   3: beatmap_thumbnail
--   4: beatmap_thumbnail_large
--   5: beatmap_audio_preview

CREATE TABLE IF NOT EXISTS resource_mirrors
(
    url character varying NOT NULL PRIMARY KEY,
    type int NOT NULL,
    server int, -- 0: ppy 1: private
    priority int DEFAULT 0
);

INSERT INTO resource_mirrors (url, type, server, priority)
VALUES ('https://osu.direct/api/osu/{}', 2, 0, 0),
       ('https://old.ppy.sh/osu/{}', 2, 0, 1),
       ('https://b.ppy.sh/thumb/{}l.jpg', 4, 0, 0),
       ('https://b.ppy.sh/thumb/{}.jpg', 3, 0, 0),
       ('https://b.ppy.sh/preview/{}.mp3', 5, 0, 0),
       ('https://api.nerinyan.moe/d/{}?noVideo=true', 1, 0, 2),
       ('https://api.nerinyan.moe/d/{}', 0, 0, 2),
       ('https://osu.direct/api/d/{}?noVideo=', 1, 0, 1),
       ('https://osu.direct/api/d/{}', 0, 0, 1),
       ('https://catboy.best/d/{}n', 1, 0, 0),
       ('https://catboy.best/d/{}', 0, 0, 0),
       ('http://keel/resources/osz/{}', 0, 1, 0),
       ('http://keel/resources/osz/{}?no_video=true', 1, 1, 0),
       ('http://keel/resources/osu/{}', 2, 1, 0),
       ('http://keel/resources/mp3/{}', 5, 1, 0),
       ('http://keel/resources/mt/{}/small', 3, 1, 0),
       ('http://keel/resources/mt/{}', 4, 1, 0);

-- Change beatmap(set) ID offset for beatmap submission
ALTER SEQUENCE beatmapsets_id_seq RESTART WITH 1000000000;
ALTER SEQUENCE beatmaps_id_seq RESTART WITH 1000000000;

CREATE INDEX IF NOT EXISTS beatmapsets_id_idx ON beatmapsets (id);
CREATE INDEX IF NOT EXISTS beatmaps_filename_idx ON beatmaps (filename);
CREATE INDEX IF NOT EXISTS beatmaps_md5_idx ON beatmaps (md5);
CREATE INDEX IF NOT EXISTS beatmaps_id_idx ON beatmaps (id);