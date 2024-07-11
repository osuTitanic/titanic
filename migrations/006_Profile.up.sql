CREATE TABLE IF NOT EXISTS profile_badges
(
    id serial NOT NULL PRIMARY KEY,
	user_id serial NOT NULL REFERENCES users (id),
	created timestamp without time zone NOT NULL DEFAULT now(),
	badge_icon character varying NOT NULL,
	badge_description character varying,
	badge_url character varying
);

CREATE TABLE IF NOT EXISTS profile_activity
(
	id serial NOT NULL PRIMARY KEY,
	user_id int NOT NULL REFERENCES users (id),
    mode smallint NOT NULL,
	"time" timestamp without time zone NOT NULL DEFAULT now(),
	activity_text character varying(256) NOT NULL,
	activity_args character varying(256),
	activity_links character varying(256)
);

CREATE TABLE IF NOT EXISTS profile_rank_history
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

CREATE TABLE IF NOT EXISTS profile_play_history
(
    user_id int NOT NULL REFERENCES users (id),
    mode smallint NOT NULL,
    year int NOT NULL,
    month int NOT NULL,
    plays int NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, year, month, mode)
);

CREATE TABLE IF NOT EXISTS profile_replay_history
(
    user_id int NOT NULL REFERENCES users (id),
    mode smallint NOT NULL,
    year int NOT NULL,
    month int NOT NULL,
    replay_views int NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, year, month, mode)
);

CREATE TABLE IF NOT EXISTS name_history
(
	id serial NOT NULL PRIMARY KEY,
	user_id int NOT NULL REFERENCES users (id),
	changed_at timestamp without time zone NOT NULL DEFAULT now(),
	name character varying NOT NULL
);