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
    grade character varying(2) NOT NULL DEFAULT 'N',
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE beatmap_nominations
(
    user_id int NOT NULL REFERENCES users (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    "time" timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, set_id)
);

CREATE TABLE beatmap_modding
(
    id serial NOT NULL PRIMARY KEY,
    target_id int NOT NULL REFERENCES users (id),
    sender_id int NOT NULL REFERENCES users (id),
    set_id int NOT NULL REFERENCES beatmapsets (id),
    post_id int NOT NULL REFERENCES forum_posts (id),
    amount int NOT NULL DEFAULT 0,
    "time" timestamp without time zone NOT NULL DEFAULT now()
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

CREATE TABLE favourites
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