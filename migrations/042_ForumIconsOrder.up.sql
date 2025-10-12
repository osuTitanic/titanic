ALTER TABLE forum_icons ADD COLUMN "order" smallint NOT NULL DEFAULT 0;

UPDATE forum_icons SET "order" = 1 WHERE name = 'star';
UPDATE forum_icons SET "order" = 2 WHERE name = 'bubble';
UPDATE forum_icons SET "order" = 3 WHERE name = 'bubblepop';
UPDATE forum_icons SET "order" = 4 WHERE name = 'heart';
UPDATE forum_icons SET "order" = 5 WHERE name = 'heartpop';
UPDATE forum_icons SET "order" = 6 WHERE name = 'fire';
UPDATE forum_icons SET "order" = 7 WHERE name = 'radioactive';
UPDATE forum_icons SET "order" = 8 WHERE name = 'question';
UPDATE forum_icons SET "order" = 9 WHERE name = 'alert';
UPDATE forum_icons SET "order" = 10 WHERE name = 'info';
UPDATE forum_icons SET "order" = 11 WHERE name = 'osu';
UPDATE forum_icons SET "order" = 12 WHERE name = 'taiko';
UPDATE forum_icons SET "order" = 13 WHERE name = 'ctb';
UPDATE forum_icons SET "order" = 14 WHERE name = 'mania';
