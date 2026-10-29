ALTER TABLE messages ADD COLUMN sender character varying(32);

UPDATE messages SET sender = users.name
FROM users WHERE users.id = messages.sender_id;

ALTER TABLE messages ALTER COLUMN sender SET NOT NULL;
CREATE INDEX idx_messages_sender_name_id_desc ON messages (sender, id DESC);
