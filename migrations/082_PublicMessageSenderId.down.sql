DROP INDEX idx_messages_sender_id_desc;
ALTER INDEX idx_messages_sender_name_id_desc RENAME TO idx_messages_sender_id_desc;
ALTER TABLE messages DROP COLUMN sender_id;
