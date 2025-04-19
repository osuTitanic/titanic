-- Extended metrics for beatmaps to make eyup star calculations possible
ALTER TABLE beatmaps ADD COLUMN diff_eyup real NOT NULL DEFAULT 0;
ALTER TABLE beatmaps ADD COLUMN count_normal int NOT NULL DEFAULT 0;
ALTER TABLE beatmaps ADD COLUMN count_slider int NOT NULL DEFAULT 0;
ALTER TABLE beatmaps ADD COLUMN count_spinner int NOT NULL DEFAULT 0;
ALTER TABLE beatmaps ADD COLUMN drain_length int NOT NULL DEFAULT 0;