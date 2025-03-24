CREATE TABLE IF NOT EXISTS quizzes
(
    id        CHAR(36) PRIMARY KEY,
    course_id CHAR(36)                                     NOT NULL,
    number    INT                                          NOT NULL,
    point     INT                                          NOT NULL,
    quiz_type ENUM ('choice', 'ordering', 'missing_words') NOT NULL,
    title     TEXT                                         NOT NULL,
    lesson    TEXT                                         NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE
);
