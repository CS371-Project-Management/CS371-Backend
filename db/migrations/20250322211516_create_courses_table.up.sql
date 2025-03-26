CREATE TABLE IF NOT EXISTS courses
(
    id               CHAR(36) PRIMARY KEY,
    class_id         CHAR(36)                        NOT NULL,
    title            VARCHAR(255)                    NOT NULL,
    description      TEXT                            NOT NULL,
    difficulty_level ENUM ('easy', 'medium', 'hard') NOT NULL,
    number           INT                             NOT NULL,
    FOREIGN KEY (class_id) REFERENCES classes (id) ON DELETE CASCADE
);