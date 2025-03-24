CREATE TABLE IF NOT EXISTS user_classes
(
    id       CHAR(36) PRIMARY KEY,
    user_id  CHAR(36) NOT NULL,
    class_id CHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (class_id) REFERENCES classes (id)
);