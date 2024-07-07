CREATE extension IF NOT EXISTS pgcrypto;

CREATE TABLE users
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

CREATE TABLE forums
(
    id serial NOT NULL PRIMARY KEY,
    parent_id int REFERENCES forums (id) DEFAULT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    name character varying(32) NOT NULL,
    description character varying(255) NOT NULL DEFAULT '',
    allow_icons boolean NOT NULL DEFAULT true,
    hidden boolean NOT NULL DEFAULT false
);

CREATE TABLE forum_icons
(
    id serial NOT NULL PRIMARY KEY,
    name character varying(32) NOT NULL,
    location character varying(255) NOT NULL
);

CREATE TABLE forum_topics
(
    id serial NOT NULL PRIMARY KEY,
    forum_id int NOT NULL REFERENCES forums (id),
    creator_id int NOT NULL REFERENCES users (id),
    title character varying(255) NOT NULL,
    status_text character varying(255) DEFAULT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    last_post_at timestamp without time zone NOT NULL DEFAULT now(),
    locked_at timestamp without time zone DEFAULT NULL,
    views int NOT NULL DEFAULT 0,
    icon smallint REFERENCES forum_icons (id) DEFAULT NULL,
    can_change_icon boolean NOT NULL DEFAULT true,
    can_star boolean NOT NULL DEFAULT false,
    announcement boolean NOT NULL DEFAULT false,
    hidden boolean NOT NULL DEFAULT false,
    pinned boolean NOT NULL DEFAULT false
);

CREATE TABLE forum_stars
(
    topic_id int NOT NULL REFERENCES forum_topics (id),
    user_id int NOT NULL REFERENCES users (id),
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE forum_posts
(
    id bigserial NOT NULL PRIMARY KEY,
    topic_id int NOT NULL REFERENCES forum_topics (id),
    forum_id int NOT NULL REFERENCES forums (id),
    user_id int NOT NULL REFERENCES users (id),
    content text NOT NULL,
    icon_id smallint REFERENCES forum_icons (id) DEFAULT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    edit_time timestamp without time zone NOT NULL DEFAULT now(),
    edit_count int NOT NULL DEFAULT 0,
    edit_locked boolean NOT NULL DEFAULT false,
    hidden boolean NOT NULL DEFAULT false,
    draft boolean NOT NULL DEFAULT false,
    deleted boolean NOT NULL DEFAULT false
);

CREATE TABLE forum_reports
(
    post_id int NOT NULL REFERENCES forum_posts (id),
    user_id int NOT NULL REFERENCES users (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    resolved_at timestamp without time zone DEFAULT NULL,
    reason character varying(255) NOT NULL,
    PRIMARY KEY (post_id, user_id)
);

CREATE TABLE forum_bookmarks
(
    user_id int NOT NULL REFERENCES users (id),
    topic_id int NOT NULL REFERENCES forum_topics (id),
    PRIMARY KEY (user_id, topic_id)
);

CREATE TABLE forum_subscribers
(
    user_id int NOT NULL REFERENCES users (id),
    topic_id int NOT NULL REFERENCES forum_topics (id),
    PRIMARY KEY (user_id, topic_id)
);

CREATE TABLE beatmapsets
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

CREATE TABLE benchmarks (
    id serial NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    smoothness real NOT NULL,
    framerate int NOT NULL,
    score int NOT NULL,
    grade character varying(2) NOT NULL DEFAULT 'N'
);

CREATE TABLE beatmap_nominations
(
    user_id int NOT NULL REFERENCES users (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    time timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, set_id)
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
    message character varying(512) NOT NULL,
    "time" timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE stats
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

CREATE TABLE screenshots
(
    id bigserial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    hidden boolean NOT NULL DEFAULT false
);

CREATE TABLE scores
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

-- Status:
-- 0: Friend
-- 1: Blocked (not implemented)
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
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, set_id)
);

CREATE TABLE achievements
(
    user_id int NOT NULL REFERENCES users (id),
    name character varying NOT NULL,
    category character varying NOT NULL,
    filename character varying NOT NULL,
    unlocked_at timestamp without time zone NOT NULL DEFAULT now(),
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
	color character varying(8),
    CONSTRAINT comments_pkey PRIMARY KEY (id)
);

CREATE TABLE profile_badges
(
    id serial NOT NULL PRIMARY KEY,
	user_id serial NOT NULL REFERENCES users (id),
	created timestamp without time zone NOT NULL DEFAULT now(),
	badge_icon character varying NOT NULL,
	badge_description character varying,
	badge_url character varying
);

CREATE TABLE profile_activity
(
	id serial NOT NULL PRIMARY KEY,
	user_id int NOT NULL REFERENCES users (id),
    mode smallint NOT NULL,
	"time" timestamp without time zone NOT NULL DEFAULT now(),
	activity_text character varying(256) NOT NULL,
	activity_args character varying(256),
	activity_links character varying(256)
);

CREATE TABLE profile_rank_history
(
   	user_id int NOT NULL REFERENCES users (id),
	"time" timestamp without time zone NOT NULL DEFAULT now(),
    mode smallint NOT NULL,
    rscore bigint NOT NULL,
    pp real NOT NULL,
    ppv1 real NOT NULL,
    pp_vn real NOT NULL,
    pp_rx real NOT NULL,
    pp_ap real NOT NULL,
    global_rank int NOT NULL,
    country_rank int NOT NULL,
    score_rank int NOT NULL,
    ppv1_rank int NOT NULL,
    pp_vn_rank int NOT NULL,
    pp_rx_rank int NOT NULL,
    pp_ap_rank int NOT NULL,
	PRIMARY KEY (user_id, "time")
);

CREATE TABLE profile_play_history
(
    user_id int NOT NULL REFERENCES users (id),
    mode smallint NOT NULL,
    year int NOT NULL,
    month int NOT NULL,
    plays int NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, year, month, mode)
);

CREATE TABLE profile_replay_history
(
    user_id int NOT NULL REFERENCES users (id),
    mode smallint NOT NULL,
    year int NOT NULL,
    month int NOT NULL,
    replay_views int NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, year, month, mode)
);

CREATE TABLE name_history
(
	id serial NOT NULL PRIMARY KEY,
	user_id int NOT NULL REFERENCES users (id),
	changed_at timestamp without time zone NOT NULL DEFAULT now(),
	name character varying NOT NULL
);

CREATE TABLE clients
(
    user_id int NOT NULL REFERENCES users (id),
    executable character(32) NOT NULL,
    adapters character(32) NOT NULL,
    unique_id character(32) NOT NULL,
    disk_signature character(32) NOT NULL,
    banned boolean NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, executable, adapters, unique_id, disk_signature)
);

-- Table for verified hardware ids, to
-- bypass multiaccounting checks.
-- Types:
-- 0: Adapters
-- 1: Unique Id
-- 2: Disk Signature
CREATE TABLE clients_verified
(
    "type" smallint NOT NULL,
    "hash" character(32) NOT NULL,
    PRIMARY KEY ("type", "hash")
);

CREATE TABLE logins
(
    user_id int NOT NULL REFERENCES users (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    ip character varying(45) NOT NULL,
    osu_version character varying(25) NOT NULL,
    PRIMARY KEY (user_id, "time")
);

CREATE TABLE reports
(
    id serial NOT NULL PRIMARY KEY,
    target_id int NOT NULL REFERENCES users (id),
    sender_id int NOT NULL REFERENCES users (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    reason character varying(255),
    resolved boolean NOT NULL DEFAULT false
);

CREATE TABLE infringements
(
    id serial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    action smallint NOT NULL DEFAULT 0, -- 0: ban 1: mute
    length timestamp without time zone,
    is_permanent boolean NOT NULL DEFAULT false,
    description character varying(255)
);

CREATE TABLE mp_matches
(
    id bigserial NOT NULL PRIMARY KEY,
    bancho_id smallint NOT NULL,
    name character varying(50) NOT NULL,
    creator_id int NOT NULL REFERENCES users (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    ended_at timestamp without time zone
);

CREATE TABLE mp_events
(
    match_id int NOT NULL REFERENCES mp_matches (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    type smallint NOT NULL,
    data jsonb NOT NULL,
    PRIMARY KEY (match_id, "time")
);

CREATE TABLE user_count
(
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    count int NOT NULL DEFAULT 0,
    PRIMARY KEY ("time")
);

CREATE TABLE verifications
(
    id serial NOT NULL PRIMARY KEY,
    token character varying(32) NOT NULL,
    user_id int NOT NULL REFERENCES users (id),
    sent_at timestamp without time zone NOT NULL DEFAULT now(),
    type smallint NOT NULL -- 0: activation - 1: password reset
);

CREATE TABLE groups
(
    id serial NOT NULL PRIMARY KEY,
    bancho_permissions smallint DEFAULT 0,
    name character varying(45) NOT NULL,
    short_name character varying(8) NOT NULL,
    description text,
    color character varying(8) NOT NULL,
    hidden boolean NOT NULL DEFAULT false
);

CREATE TABLE groups_entries
(
    group_id int NOT NULL REFERENCES groups (id),
    user_id int NOT NULL REFERENCES users (id),
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE notifications
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

CREATE TABLE releases
(
    name character varying(64) NOT NULL PRIMARY KEY,
    version int NOT NULL,
    description text NOT NULL,
    known_bugs text,
    supported boolean NOT NULL DEFAULT true,
    recommended boolean NOT NULL DEFAULT false,
    preview boolean NOT NULL DEFAULT false,
    downloads varchar[] NOT NULL DEFAULT '{}',
    hashes jsonb NOT NULL DEFAULT '[]',
    screenshots jsonb NOT NULL DEFAULT '[]',
    actions jsonb NOT NULL DEFAULT '[]',
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

-- resource type:
--   0: osz
--   1: osz_novideo
--   2: beatmap_file
--   3: beatmap_thumbnail
--   4: beatmap_thumbnail_large
--   5: beatmap_audio_preview

CREATE TABLE resource_mirrors
(
    url character varying NOT NULL PRIMARY KEY,
    type int NOT NULL,
    server int, -- 0: ppy 1: private
    priority int DEFAULT 0
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

INSERT INTO channels (name, topic, read_permissions, write_permissions)
VALUES ('#osu', 'General discussion.', 1, 1),
       ('#announce', 'Public announcements.', 1, 8),
       ('#lobby', 'Multiplayer lobby discussion room.', 1, 1),
       ('#admin', 'General discussion for administrators.', 16, 16);

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

INSERT INTO clients_verified ("type", "hash")
VALUES (0, 'b4ec3c4334a0249dae95c284ec5983df'), -- "runningunderwine"
       (0, '74be16979710d4c4e7c6647856088456'), -- ""
       (0, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (1, 'ad921d60486366258809553a3db49a4a'), -- "unknown"
       (1, '74be16979710d4c4e7c6647856088456'), -- ""
       (1, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (2, 'ad921d60486366258809553a3db49a4a'), -- "unknown"
       (2, 'dcfcd07e645d245babe887e5e2daa016'), -- "0"
       (2, '28c8edde3d61a0411511d3b1866f0636'), -- "1"
       (2, '74be16979710d4c4e7c6647856088456'), -- ""
       (2, 'd41d8cd98f00b204e9800998ecf8427e'), -- ""
       (2, 'd1c651c36f499849f1c9a5843567e686'); -- toshiba hdd

INSERT INTO user_count (count)
VALUES (0);

INSERT INTO forums (id, name)
VALUES (1, 'Titanic!'),
       (2, 'Beatmaps');

INSERT INTO forums (id, name, description, parent_id, allow_icons)
VALUES (3, 'Development', 'Discuss the future of this project.', 1, true),
       (4, 'Gameplay & Rankings', 'Show off your scores to the world and discuss them.', 1, true),
       (5, 'Skinning', 'Discuss & share skins and other customizations.', 1, true),
       (6, 'Feature Requests', 'Suggest what you would like to see in this project.', 1, true),
       (7, 'Support', 'Need help? You will find it here.', 1, true),
       (8, 'Ranked/Approved Beatmaps', 'New approved beatmaps will show up in here.', 2, false),
       (9, 'Pending Beatmaps', 'New pending beatmaps that are waiting for approval.', 2, false),
       (10, 'Work In Progress/Help Wanted', 'Work-in-progress beatmaps that may need support/help.', 2, false),
       (11, 'Map Requests', 'Request beatmaps from the official servers.', 2, false),
       (12, 'Beatmap Graveyard', 'Beatmaps that havent been active for 4 weeks or more will be moved here.', 2, false),
       (13, 'Mapping Discussion', 'Share the art of mapping with others.', 2, false);

INSERT INTO forums (id, name, parent_id, allow_icons)
VALUES (14, 'Taiko', 4, false),
       (15, 'Catch the Beat', 4, false),
       (16, 'osu!mania', 4, false),
       (17, 'Completed Skins', 5, false),
       (18, 'Bug Reports', 3, false);

INSERT INTO forums (id, name)
VALUES (19, 'Other');

INSERT INTO forums (id, name, description, parent_id, allow_icons)
VALUES (20, 'General Discussion', 'The place where you dont post crap.', 19, false),
       (21, 'Off-Topic', 'The perfect place for brainrot.', 19, false),
       (22, 'Introductions', 'Introduce yourself to other passengers.', 19, false),
       (23, 'Client Modding', 'Discover a new world of osu!.', 19, false),
       (24, 'Video Games', 'Discuss any non-osu! games in here.', 19, false),
       (25, 'Art', 'Share your artistic masterpieces, find new avatars and more.', 19, false);

INSERT INTO forum_icons (name, location)
VALUES ('heart', '/images/icons/forum/heart.gif'),
       ('heartpop', '/images/icons/forum/heartpop.gif'),
       ('bubble', '/images/icons/forum/thinking.gif'),
       ('bubblepop', '/images/icons/forum/bubblepop.png'),
       ('fire', '/images/icons/forum/fire.gif'),
       ('star', '/images/icons/forum/star.gif'),
       ('radioactive', '/images/icons/forum/radioactive.gif'),
       ('alert', '/images/icons/forum/alert.gif'),
       ('info', '/images/icons/forum/info.gif'),
       ('question', '/images/icons/forum/question.gif'),
       ('osu', '/images/icons/forum/osu.gif'),
       ('taiko', '/images/icons/forum/taiko.gif'),
       ('ctb', '/images/icons/forum/ctb.gif'),
       ('mania', '/images/icons/forum/mania.gif');

INSERT INTO resource_mirrors (url, type, server, priority)
VALUES ('https://api.osu.direct/osu/{}', 2, 0, 0),
       ('https://old.ppy.sh/osu/{}', 2, 0, 1),
       ('https://b.ppy.sh/thumb/{}l.jpg', 4, 0, 0),
       ('https://b.ppy.sh/thumb/{}.jpg', 3, 0, 0),
       ('https://b.ppy.sh/preview/{}.mp3', 5, 0, 0),
       ('https://api.nerinyan.moe/d/{}?noVideo=true', 1, 0, 1),
       ('https://api.nerinyan.moe/d/{}', 0, 0, 1),
       ('https://api.osu.direct/d/{}?noVideo=', 1, 0, 0),
       ('https://api.osu.direct/d/{}', 0, 0, 0),
       ('/api/beatmaps/osz/{}', 0, 1, 0),
       ('/api/beatmaps/osz/{}?noVideo=true', 1, 1, 0),
       ('/api/beatmaps/osu/{}', 2, 1, 0),
       ('/api/beatmaps/mp3/{}', 5, 1, 0),
       ('/api/beatmaps/mt/{}', 3, 1, 0),
       ('/api/beatmaps/mt/{}?large=true', 4, 1, 0);

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

-- Change beatmap(set) ID offset for beatmap submission
ALTER SEQUENCE beatmapsets_id_seq RESTART WITH 1000000000;
ALTER SEQUENCE beatmaps_id_seq RESTART WITH 1000000000;
