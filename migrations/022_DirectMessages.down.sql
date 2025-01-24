-- Revert migration
INSERT INTO messages (sender, target, message, "time")
SELECT
    sender_user.name AS sender,
    target_user.name AS target,
    dm.message,
    dm."time"
FROM
    direct_messages AS dm
JOIN
    users AS sender_user ON sender_user.id = dm.sender_id
JOIN    
    users AS target_user ON target_user.id = dm.target_id;

DELETE FROM direct_messages;
DROP TABLE direct_messages;