-- Groups & forums use hardcoded values for their IDs, so we
-- need to reset the sequences to avoid conflicts when creating
-- new groups or forums

SELECT setval(
    'groups_id_seq',
    GREATEST(
        (SELECT last_value FROM groups_id_seq),
        (SELECT COALESCE(MAX(id), 1) FROM groups)
    )
);

SELECT setval(
    'forums_id_seq',
    GREATEST(
        (SELECT last_value FROM forums_id_seq),
        (SELECT COALESCE(MAX(id), 1) FROM forums)
    )
);
