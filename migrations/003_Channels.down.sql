DELETE FROM channels
WHERE name IN (
    '#osu',
    '#announce',
    '#lobby',
    '#admin'
);