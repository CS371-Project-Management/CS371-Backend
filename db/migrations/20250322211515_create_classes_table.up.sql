CREATE TABLE IF NOT EXISTS classes
(
    id            CHAR(36) PRIMARY KEY,
    user_id       CHAR(36)     NOT NULL,
    invite_code   VARCHAR(255) NOT NULL,
    title         VARCHAR(255) NOT NULL,
    description   TEXT         NOT NULL,
    accessibility BOOLEAN      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);