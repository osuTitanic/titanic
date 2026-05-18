-- Migration removes wildcard in favor of explicit permissions for kudosu management
-- This is due to a new `forum.kudosu.force_reward` permission that only certain users should have, not the entire group

DELETE FROM group_permissions
WHERE group_id = 3
    AND permission = 'forum.kudosu.*'
    AND rejected = false;

INSERT INTO group_permissions (group_id, permission, rejected)
SELECT group_id, permission, false
FROM (VALUES
    (3, 'forum.kudosu.reward'),
    (3, 'forum.kudosu.revoke'),
    (3, 'forum.kudosu.reset')
) AS new_permissions(group_id, permission)
WHERE NOT EXISTS (
    SELECT 1
    FROM group_permissions existing
    WHERE existing.group_id = new_permissions.group_id
      AND existing.permission = new_permissions.permission
      AND existing.rejected = false
);
