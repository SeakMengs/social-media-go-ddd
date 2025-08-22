CREATE TABLE follows (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    follower_id CHAR(36) NOT NULL,
    followee_id CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (follower_id, followee_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE
);
