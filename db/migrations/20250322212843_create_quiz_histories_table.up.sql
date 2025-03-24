CREATE TABLE IF NOT EXISTS quiz_histories
(
    id        CHAR(36) PRIMARY KEY,
    user_id   CHAR(36)                                     NOT NULL,
    quiz_id   CHAR(36)                                     NOT NULL,
    point     INT                                          NOT NULL,
    quiz_type ENUM ('choice', 'ordering', 'missing_words') NOT NULL,
    note      TEXT,
    passed    BOOLEAN                                      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id)
);