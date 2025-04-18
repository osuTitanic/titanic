DELETE FROM user_permissions;
DELETE FROM group_permissions;

DROP INDEX IF EXISTS idx_user_permissions_user_id;
DROP INDEX IF EXISTS idx_group_permissions_group_id;
DROP INDEX IF EXISTS idx_user_permissions_user_id_rejected;
DROP INDEX IF EXISTS idx_group_permissions_group_id_rejected;

DROP TABLE user_permissions;
DROP TABLE group_permissions;