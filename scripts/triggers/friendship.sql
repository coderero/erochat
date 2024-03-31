CREATE TRIGGER check_column_duplication BEFORE INSERT ON friendships FOR EACH ROW BEGIN IF EXISTS (
    SELECT
        1
    FROM
        friendships
    WHERE
        (
            user_id = NEW.user_id
            AND friend_id = NEW.friend_id
        )
        OR (
            user_id = NEW.friend_id
            AND friend_id = NEW.user_id
        )
) THEN SIGNAL SQLSTATE '45000'
SET
    MESSAGE_TEXT = 'exception: self request';

END IF;

END;

CREATE TRIGGER check_request_duplication BEFORE INSERT ON friendships FOR EACH ROW BEGIN DECLARE f_count INT;

SELECT
    COUNT(*) INTO f_count
FROM
    friendships
WHERE
    (
        user1 = NEW.user1
        and user2 = NEW.user2
    )
    OR (
        user1 = NEW.user2
        and user2 = NEW.user1
    );

IF f_count > 0 THEN SIGNAL SQLSTATE '45000'
SET
    MESSAGE_TEXT = 'exception: duplicate request';

END IF;

END;