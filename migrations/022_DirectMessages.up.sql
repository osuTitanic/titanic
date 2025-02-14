CREATE TABLE IF NOT EXISTS direct_messages
(
    id bigserial NOT NULL PRIMARY KEY,
    sender_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    target_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    message text NOT NULL,
    "time" timestamp without time zone NOT NULL DEFAULT now()
);

CREATE INDEX dms_sender_id_idx ON direct_messages (sender_id, id DESC);
CREATE INDEX dms_target_id_idx ON direct_messages (target_id, id DESC);

-- Migrate previous messages to the new table
INSERT INTO direct_messages (sender_id, target_id, message, "time")
SELECT
    sender_user.id AS sender_id,
    target_user.id AS target_id,
    msg.message,
    msg."time"
FROM
    messages AS msg
JOIN
    users AS sender_user ON sender_user.name = msg.sender
JOIN
    users AS target_user ON target_user.name = msg.target
WHERE
    msg.target NOT LIKE '#%';

-- Delete direct messages from old table
DELETE FROM messages WHERE target NOT LIKE '#%';