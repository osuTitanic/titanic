CREATE TABLE IF NOT EXISTS logins
(
    user_id int NOT NULL REFERENCES users (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    ip character varying(45) NOT NULL,
    osu_version character varying(25) NOT NULL,
    PRIMARY KEY (user_id, "time")
);

CREATE TABLE IF NOT EXISTS channels
(
    name character varying(32) NOT NULL PRIMARY KEY,
    topic character varying(128) NOT NULL,
    read_permissions int NOT NULL DEFAULT 1,
    write_permissions int NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS messages
(
    id bigserial NOT NULL PRIMARY KEY,
    sender character varying(32) NOT NULL,
    target character varying(32) NOT NULL,
    message character varying(512) NOT NULL,
    "time" timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS logs
(
    id bigserial NOT NULL PRIMARY KEY,
    level character varying(12) NOT NULL,
    type character varying(250) NOT NULL,
    message character varying NOT NULL,
    "time" timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS reports
(
    id serial NOT NULL PRIMARY KEY,
    target_id int NOT NULL REFERENCES users (id),
    sender_id int NOT NULL REFERENCES users (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    reason character varying(255),
    resolved boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS mp_matches
(
    id bigserial NOT NULL PRIMARY KEY,
    bancho_id smallint NOT NULL,
    name character varying(50) NOT NULL,
    creator_id int NOT NULL REFERENCES users (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    ended_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS mp_events
(
    match_id int NOT NULL REFERENCES mp_matches (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    type smallint NOT NULL,
    data jsonb NOT NULL,
    PRIMARY KEY (match_id, "time")
);

INSERT INTO channels (name, topic, read_permissions, write_permissions)
VALUES ('#osu', 'General discussion.', 1, 1),
       ('#announce', 'Public announcements.', 1, 8),
       ('#lobby', 'Multiplayer lobby discussion room.', 1, 1),
       ('#admin', 'General discussion for administrators.', 16, 16);