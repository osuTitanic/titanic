CREATE TABLE IF NOT EXISTS profile_stamps
(
    id serial NOT NULL PRIMARY KEY,
	user_id serial NOT NULL REFERENCES users (id),
	created timestamp without time zone NOT NULL DEFAULT now(),
	icon character varying NOT NULL,
	description character varying,
	url character varying
);

CREATE INDEX IF NOT EXISTS idx_profile_stamps_user_id ON profile_stamps (user_id);
