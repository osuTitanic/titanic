DELETE FROM group_permissions
WHERE rejected = false
AND group_id = 3
AND permission IN ('forum.topics.create_beatmap', 'forum.topics.link_beatmapset');
