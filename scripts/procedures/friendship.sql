CREATE PROCEDURE get_friend (
    IN user_uid VARCHAR(36),
    IN friend_uid VARCHAR(36)
) BEGIN
SELECT
    f.uid AS rid,
    IF (u1.uid = user_uid, u2.uid, u1.uid) AS uid,
    IF (u1.uid = user_uid, u2.username, u1.username) AS username,
    IF (u1.uid = user_uid, p2.first_name, p1.first_name) AS first_name,
    IF (u1.uid = user_uid, p2.last_name, p1.last_name) AS last_name,
    IF (u1.uid = user_uid, p2.bio, p1.bio) AS bio,
    IF (u1.uid = user_uid, p2.avatar, p1.avatar) AS avatar,
    f.accepted_at
FROM
    friendships f
    JOIN users u1 ON u1.uid = f.user1
    JOIN users u2 ON u2.uid = f.user2
    LEFT JOIN profiles p1 ON p1.uid = f.user1
    LEFT JOIN profiles p2 ON p2.uid = f.user2
WHERE
    f.accepted = TRUE
    AND (
        u1.deleted_at IS NULL
        OR u2.deleted_at IS NULL
    )
    AND (
        user1 = user_uid
        OR user2 = user_uid
    )
    AND (
        user1 = friend_uid
        OR user2 = friend_uid
    );

END;

CREATE PROCEDURE get_friends_or_requests (IN user_uid VARCHAR(36), IN request BOOLEAN) BEGIN
SELECT
    f.uid AS rid,
    IF (u1.uid = user_uid, u2.uid, u1.uid) AS uid,
    IF (u1.uid = user_uid, u2.username, u1.username) AS username,
    IF (u1.uid = user_uid, p2.first_name, p1.first_name) AS first_name,
    IF (u1.uid = user_uid, p2.last_name, p1.last_name) AS last_name,
    IF (u1.uid = user_uid, p2.bio, p1.bio) AS bio,
    IF (u1.uid = user_uid, p2.avatar, p1.avatar) AS avatar,
    f.accepted_at
FROM
    friendships f
    JOIN users u1 ON u1.uid = f.user1
    JOIN users u2 ON u2.uid = f.user2
    LEFT JOIN profiles p1 ON p1.uid = f.user1
    LEFT JOIN profiles p2 ON p2.uid = f.user2
WHERE
    f.accepted = (NOT request)
    AND (
        u1.deleted_at IS NULL
        OR u2.deleted_at IS NULL
    )
    AND (
        user1 = user_uid
        OR user2 = user_uid
    );

END;

CREATE PROCEDURE get_friend_request (
    IN user_uid VARCHAR(36),
    IN request_id VARCHAR(36)
) BEGIN
SELECT
    f.uid AS rid,
    IF (u1.uid = user_uid, u2.uid, u1.uid) AS uid,
    IF (u1.uid = user_uid, u2.username, u1.username) AS username,
    IF (u1.uid = user_uid, p2.first_name, p1.first_name) AS first_name,
    IF (u1.uid = user_uid, p2.last_name, p1.last_name) AS last_name,
    IF (u1.uid = user_uid, p2.bio, p1.bio) AS bio,
    IF (u1.uid = user_uid, p2.avatar, p1.avatar) AS avatar,
    f.accepted_at
FROM
    friendships f
    JOIN users u1 ON u1.uid = f.user1
    JOIN users u2 ON u2.uid = f.user2
    LEFT JOIN profiles p1 ON p1.uid = f.user1
    LEFT JOIN profiles p2 ON p2.uid = f.user2
WHERE
    f.uid = request_id
    AND (
        u1.deleted_at IS NULL
        OR u2.deleted_at IS NULL
    )
    AND f.accepted_at IS NULL;

END;

CREATE PROCEDURE get_friends_status (IN user_uid VARCHAR(36)) BEGIN
SELECT
    f.uid AS rid,
    IF (f.user1 = user_uid, u2.uid, u1.uid) AS uid,
    s.uid AS status_id,
    IF (f.user1 = user_uid, u2.username, u1.username) AS username,
    IF (f.user1 = user_uid, p2.first_name, p1.first_name) AS first_name,
    IF (f.user1 = user_uid, p2.last_name, p1.last_name) AS last_name,
    IF (f.user1 = user_uid, p2.avatar, p1.avatar) AS avatar,
    s.title,
    s.resource_uri,
    s.resource_thumbnail
FROM
    friendships f
    JOIN users u1 ON f.user1 = u1.uid
    JOIN users u2 ON f.user2 = u2.uid
    JOIN profiles p1 ON u1.id = p1.user_id
    JOIN profiles p2 ON u2.id = p2.user_id
    LEFT JOIN status s ON (
        (
            s.user_uid = u1.uid
            AND f.user1 <> user_uid
        )
        OR (
            s.user_uid = u2.uid
            AND f.user2 <> user_uid
        )
    )
WHERE
    f.accepted = TRUE
    AND (
        user_uid = f.user1
        OR user_uid = f.user2
    )
    AND (
        u1.deleted_at IS NULL
        AND u2.deleted_at IS NULL
    )
    AND (
        s.created_at > DATE_SUB (NOW (), INTERVAL 1 DAY)
        AND s.deleted_at IS NULL
    )
    AND (s.uid IS NOT NULL);

END;

CREATE PROCEDURE get_friend_status (IN user_uid VARCHAR(36), IN status_id VARCHAR(36)) BEGIN
SELECT
    f.uid AS rid,
    IF (f.user1 = user_uid, u2.uid, u1.uid) AS uid,
    s.uid AS status_id,
    IF (f.user1 = user_uid, u2.username, u1.username) AS username,
    IF (f.user1 = user_uid, p2.first_name, p1.first_name) AS first_name,
    IF (f.user1 = user_uid, p2.last_name, p1.last_name) AS last_name,
    IF (f.user1 = user_uid, p2.avatar, p1.avatar) AS avatar,
    s.title,
    s.resource_uri,
    s.resource_thumbnail
FROM
    friendships f
    JOIN users u1 ON f.user1 = u1.uid
    JOIN users u2 ON f.user2 = u2.uid
    JOIN profiles p1 ON u1.id = p1.user_id
    JOIN profiles p2 ON u2.id = p2.user_id
    LEFT JOIN status s ON (
        (
            s.user_uid = u1.uid
            AND f.user1 <> user_uid
        )
        OR (
            s.user_uid = u2.uid
            AND f.user2 <> user_uid
        )
    )
WHERE
    f.accepted = TRUE
    AND (
        user_uid = f.user1
        OR user_uid = f.user2
    )
    AND (
        u1.deleted_at IS NULL
        AND u2.deleted_at IS NULL
    )
    AND (
        s.created_at > DATE_SUB (NOW (), INTERVAL 1 DAY)
        AND s.deleted_at IS NULL
    )
    AND (s.uid IS NOT NULL)
    AND (s.uid = status_id);

END;