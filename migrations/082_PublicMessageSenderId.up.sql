ALTER TABLE messages ADD COLUMN sender_id INTEGER REFERENCES users (id);

ALTER INDEX idx_messages_sender_id_desc RENAME TO idx_messages_sender_name_id_desc;
CREATE INDEX idx_messages_sender_id_desc ON messages (sender_id, id DESC);

UPDATE messages SET sender_id = 1 WHERE sender_id IS NULL;
ALTER TABLE messages ALTER COLUMN sender_id SET NOT NULL;
