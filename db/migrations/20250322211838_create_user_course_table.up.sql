CREATE TABLE IF NOT EXISTS user_courses
(
    id        CHAR(36) PRIMARY KEY,
    user_id   CHAR(36)                                                NOT NULL,
    course_id CHAR(36)                                                NOT NULL,
    status    ENUM ('in_progress', 'completed') DEFAULT 'in_progress' NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE
);