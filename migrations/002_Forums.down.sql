DELETE FROM forum_icons
WHERE name IN (
    'heart',
    'heartpop',
    'bubble',
    'bubblepop',
    'fire',
    'star',
    'radioactive',
    'alert',
    'info',
    'question',
    'osu',
    'taiko',
    'ctb',
    'mania'
);

DELETE FROM forums
WHERE id IN (
    1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
    14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25
);

DROP TABLE IF EXISTS forum_subscribers;
DROP TABLE IF EXISTS forum_bookmarks;
DROP TABLE IF EXISTS forum_reports;
DROP TABLE IF EXISTS forum_posts;
DROP TABLE IF EXISTS forum_stars;
DROP TABLE IF EXISTS forum_topics;
DROP TABLE IF EXISTS forum_icons;
DROP TABLE IF EXISTS forums;