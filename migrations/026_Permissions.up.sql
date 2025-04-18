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