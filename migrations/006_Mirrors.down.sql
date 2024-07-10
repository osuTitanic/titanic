DELETE FROM resource_mirrors
WHERE url IN (
    'https://api.osu.direct/osu/{}',
    'https://old.ppy.sh/osu/{}',
    'https://b.ppy.sh/thumb/{}l.jpg',
    'https://b.ppy.sh/thumb/{}.jpg',
    'https://b.ppy.sh/preview/{}.mp3',
    'https://api.nerinyan.moe/d/{}?noVideo=true',
    'https://api.nerinyan.moe/d/{}',
    'https://api.osu.direct/d/{}?noVideo=',
    'https://api.osu.direct/d/{}',
    '/api/beatmaps/osz/{}',
    '/api/beatmaps/osz/{}?noVideo=true',
    '/api/beatmaps/osu/{}',
    '/api/beatmaps/mp3/{}',
    '/api/beatmaps/mt/{}',
    '/api/beatmaps/mt/{}?large=true'
);