-- Change beatmap(set) ID offset for beatmap submission
ALTER SEQUENCE beatmapsets_id_seq RESTART WITH 1000000000;
ALTER SEQUENCE beatmaps_id_seq RESTART WITH 1000000000;