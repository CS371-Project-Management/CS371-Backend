CREATE TABLE IF NOT EXISTS users
(
    id         CHAR(36) PRIMARY KEY,
    email      VARCHAR(255)                NOT NULL UNIQUE,
    first_name VARCHAR(255)                NOT NULL,
    last_name  VARCHAR(255)                NOT NULL,
    password   VARCHAR(255)                NOT NULL,
    username   VARCHAR(255)                NOT NULL UNIQUE,
    role       ENUM ('student', 'teacher') NOT NULL
);
