CREATE TABLE
    users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        uid VARCHAR(36) NOT NULL UNIQUE,
        username VARCHAR(255) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP NULL,
    );

CREATE TABLE
    profiles (
        id INT AUTO_INCREMENT PRIMARY KEY,
        uid VARCHAR(36) NOT NULL UNIQUE,
        user_id INT NOT NULL UNIQUE,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
        bio TEXT NOT NULL,
        avatar VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP NULL,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (uid) REFERENCES users (uid)
    );

CREATE TABLE
    friendships (
        id INT AUTO_INCREMENT PRIMARY KEY,
        uid VARCHAR(36) NOT NULL UNIQUE,
        user1 VARCHAR(36) NOT NULL,
        user2 VARCHAR(36) NOT NULL,
        accepted BOOLEAN DEFAULT false NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        accepted_at TIMESTAMP NULL,
        FOREIGN KEY (user1) REFERENCES users (id),
        FOREIGN KEY (user2) REFERENCES users (id),
    );

CREATE TABLE
    status (
        id INT AUTO_INCREMENT PRIMARY KEY,
        uid VARCHAR(36) NOT NULL UNIQUE,
        user_uid VARCHAR(36) NOT NULL,
        resource_uri VARCHAR(255) NOT NULL UNIQUE,
        resource_thumbnail VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        title VARCHAR(255) NOT NULL,
        deleted_at TIMESTAMP NULL,
        FOREIGN KEY (user_uid) REFERENCES users (uid)
    )