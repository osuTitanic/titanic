CREATE TABLE beatmap_packs (
    id serial PRIMARY KEY,
    name character varying(255) NOT NULL,
    category character varying(255) NOT NULL,
    download_link character varying(255) NOT NULL,
    description text NOT NULL DEFAULT '',
    creator_id int NOT NULL REFERENCES users (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE beatmap_pack_entries (
    pack_id int NOT NULL REFERENCES beatmap_packs (id),
    beatmapset_id int NOT NULL REFERENCES beatmapsets (id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (pack_id, beatmapset_id)
);

CREATE INDEX beatmap_packs_creator_id ON beatmap_packs (creator_id);
CREATE INDEX beatmap_packs_category ON beatmap_packs (category);

CREATE INDEX beatmap_pack_entries_pack_id ON beatmap_pack_entries (pack_id);
CREATE INDEX beatmap_pack_entries_beatmapset_id ON beatmap_pack_entries (beatmapset_id);
