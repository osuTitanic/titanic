DELETE FROM group_permissions
WHERE group_id = 3
    AND permission IN (
        'forum.kudosu.reward',
        'forum.kudosu.revoke',
        'forum.kudosu.reset'
    )
    AND rejected = false;

INSERT INTO group_permissions (group_id, permission, rejected)
SELECT 3, 'forum.kudosu.*', false
WHERE NOT EXISTS (
    SELECT 1
    FROM group_permissions
    WHERE group_id = 3
        AND permission = 'forum.kudosu.*'
        AND rejected = false
);
