CREATE TABLE user_permissions (
    id serial PRIMARY KEY,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission varchar(255) NOT NULL,
    rejected boolean DEFAULT false,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
);

CREATE TABLE group_permissions (
    id serial PRIMARY KEY,
    group_id int NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission varchar(255) NOT NULL,
    rejected boolean DEFAULT false,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
);

CREATE INDEX idx_user_permissions_user_id ON user_permissions(user_id);
CREATE INDEX idx_group_permissions_group_id ON group_permissions(group_id);
CREATE INDEX idx_user_permissions_user_id_rejected ON user_permissions(user_id, rejected);
CREATE INDEX idx_group_permissions_group_id_rejected ON group_permissions(group_id, rejected);