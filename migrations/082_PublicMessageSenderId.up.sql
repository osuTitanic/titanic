ALTER TABLE messages ADD COLUMN sender_id INTEGER REFERENCES users (id);

ALTER INDEX idx_messages_sender_id_desc RENAME TO idx_messages_sender_name_id_desc;
CREATE INDEX idx_messages_sender_id_desc ON messages (sender_id, id DESC);
