DELETE FROM groups_entries
WHERE user_id = 1 AND group_id IN (1, 999, 1000) OR
      user_id = 2 AND group_id IN (1, 999, 1000);

DELETE FROM groups
WHERE id IN (
    1, 2, 3, 4, 5, 6, 7, 8,
    9, 997, 998, 999, 1000
);