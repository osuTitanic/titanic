INSERT INTO group_permissions (group_id, permission, rejected)
SELECT group_id, permission, false
FROM (VALUES
    (3, 'forum.topics.create_beatmap'),
    (3, 'forum.topics.link_beatmapset')
) AS new_permissions(group_id, permission)
WHERE NOT EXISTS (
    -- Check if the permission already exists for the group
    SELECT 1
    FROM group_permissions existing
    WHERE existing.group_id = new_permissions.group_id
        AND existing.permission = new_permissions.permission
        AND existing.rejected = false
);
