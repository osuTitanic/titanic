CREATE TABLE IF NOT EXISTS users
(
    id serial NOT NULL PRIMARY KEY,
    name character varying(32) NOT NULL,
    safe_name character varying(32) NOT NULL,
    email character varying(255) NOT NULL,
    pw character(60) NOT NULL, -- bcrypt
    discord_id bigint,
    bot boolean NOT NULL DEFAULT false,
    country character varying NOT NULL DEFAULT 'XX',
    silence_end timestamp without time zone,
    supporter_end timestamp without time zone,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    latest_activity timestamp without time zone NOT NULL DEFAULT now(),
    restricted boolean NOT NULL DEFAULT false,
    activated boolean NOT NULL DEFAULT false,
    preferred_mode int NOT NULL DEFAULT 0,
    playstyle int NOT NULL DEFAULT 0,
    kudosu int NOT NULL DEFAULT 0,
    irc_token character(10) NOT NULL DEFAULT encode(gen_random_bytes(5), 'hex'),
    userpage_about text,
    userpage_signature text,
    userpage_title character varying(64),
    userpage_banner character varying(255),
    userpage_website character varying(64),
    userpage_discord character varying(64),
    userpage_twitter character varying(64),
    userpage_location character varying(30),
    userpage_interests character varying(30),
    UNIQUE(name, safe_name, email, discord_id)
);

CREATE TABLE IF NOT EXISTS stats
(
    id int NOT NULL REFERENCES users (id),
    mode smallint NOT NULL, -- osu!, taiko, ctb and mania
    rank int NOT NULL DEFAULT 0,
    tscore bigint NOT NULL DEFAULT 0,
    rscore bigint NOT NULL DEFAULT 0,
    pp real NOT NULL DEFAULT 0,
    ppv1 real NOT NULL DEFAULT 0,
    pp_vn real NOT NULL DEFAULT 0,
    pp_rx real NOT NULL DEFAULT 0,
    pp_ap real NOT NULL DEFAULT 0,
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

CREATE TABLE IF NOT EXISTS groups
(
    id serial NOT NULL PRIMARY KEY,
    bancho_permissions smallint DEFAULT 0,
    name character varying(45) NOT NULL,
    short_name character varying(8) NOT NULL,
    description text,
    color character varying(8) NOT NULL,
    hidden boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS groups_entries
(
    group_id int NOT NULL REFERENCES groups (id),
    user_id int NOT NULL REFERENCES users (id),
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE IF NOT EXISTS user_count
(
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    count int NOT NULL DEFAULT 0,
    PRIMARY KEY ("time")
);

-- Status:
-- 0: Friend
-- 1: Blocked
CREATE TABLE IF NOT EXISTS relationships
(
    user_id int NOT NULL REFERENCES users (id),
    target_id int NOT NULL REFERENCES users (id),
    status smallint NOT NULL,
    PRIMARY KEY (user_id, target_id)
);

CREATE TABLE IF NOT EXISTS infringements
(
    id serial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    action smallint NOT NULL DEFAULT 0, -- 0: ban 1: mute
    length timestamp without time zone,
    is_permanent boolean NOT NULL DEFAULT false,
    description character varying(255)
);

CREATE TABLE IF NOT EXISTS verifications
(
    id serial NOT NULL PRIMARY KEY,
    token character varying(32) NOT NULL,
    user_id int NOT NULL REFERENCES users (id),
    sent_at timestamp without time zone NOT NULL DEFAULT now(),
    type smallint NOT NULL -- 0: activation - 1: password reset
);

CREATE TABLE IF NOT EXISTS notifications
(
    id bigserial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users (id),
    type smallint NOT NULL,
    header character varying(128) NOT NULL,
    content text NOT NULL,
    link character varying(255),
    read boolean NOT NULL DEFAULT false,
    "time" timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS screenshots
(
    id bigserial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    hidden boolean NOT NULL DEFAULT false
);

INSERT INTO users (name, safe_name, email, pw, country, activated, bot)
VALUES ('BanchoBot', 'banchobot', 'bot@example.com', '------------------------------------------------------------', 'OC', true, true),
       ('peppy', 'peppy', 'pe@ppy.sh', '$2b$12$W5ppLwlSEJ3rpJQRq8UcX.QA5cTm7HvsVpn6MXQHE/6OEO.Iv4DGW', 'AU', true, false);

-- Default stats for BanchoBot
INSERT INTO stats (id, mode)
VALUES (1, 0),
       (1, 1),
       (1, 2),
       (1, 3);

INSERT INTO groups (id, bancho_permissions, name, short_name, description, color, hidden)
VALUES (1, '16', 'Admins', 'ADMIN', 'Some cool people.', '#9d6b15', false),
       (2, '8', 'Developers', 'DEV', 'These people actively contribute to the titanic codebase.', '#69159d', false),
       (3, '2', 'Beatmap Approval Team', 'BAT', 'These people handle map requests from users and decide which should be approved.', '#1959b8', false),
       (4, '8', 'Global Moderator Team', 'GMT', 'These people focus on player moderation and ensuring the community is a safe place for everyone.', '#59c51b', false),
       (5, '32', 'Tournament Manager Team', 'TMT', 'These people help organize & referee tournament.', '#2117ab', false),
       (6, 4, 'Donators', 'DONOR', 'These people have donated money to this project and are helping to keep it alive. Eternal thanks to them!', '#cf9c02', false),
       (7, 0, 'Alumni', 'ALM', 'These people have voluntary contributed to this project in some way.', '#c51b54', false),
       (8, 0, 'Playtesters', 'TESTER', 'These people were one of the first to register on this server and have helped us finding bugs and giving feedback to the project.', '#31b6e3', false),
       (9, 0, 'Bots', 'BOT', '*beep boop*', '#00b061', false),
       (997, 0, 'Preview', 'PREVIEW', 'These people can access unverified builds of the game.', '#000000', true),
       (998, '1', 'Verified', 'VERIFIED', 'Verified players.', '#000000', true),
       (999, '4', 'Supporter', 'DIRECT', 'People with access to osu! direct.', '#000000', true),
       (1000, '1', 'Players', 'PLAYER', 'People who play the game.', '#000000', true);

INSERT INTO groups_entries (user_id, group_id)
VALUES (1, 1),    -- BanchoBot -> Admins
       (1, 999),  -- BanchoBot -> Supporter
       (1, 1000), -- BanchoBot -> Players
       (2, 1),    -- peppy -> Admins
       (2, 999),  -- peppy -> Supporter
       (2, 1000); -- peppy -> Players

INSERT INTO user_count (count)
VALUES (0);

CREATE INDEX IF NOT EXISTS group_id_idx ON groups_entries (group_id);
CREATE INDEX IF NOT EXISTS users_name_idx ON users (name);
CREATE INDEX IF NOT EXISTS users_id_idx ON users (id);
CREATE INDEX IF NOT EXISTS stats_id_idx ON stats (id);