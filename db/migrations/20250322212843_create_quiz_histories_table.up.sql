CREATE TABLE IF NOT EXISTS quiz_histories
(
    id              CHAR(36) PRIMARY KEY,
    user_id         CHAR(36)                                     NOT NULL,
    quiz_history_id CHAR(36)                                     NOT NULL,
    quiz_type       ENUM ('choice', 'ordering', 'missing_words') NOT NULL,
    note            TEXT,
    passed          BOOLEAN                                      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (quiz_history_id) REFERENCES quizzes (id) ON DELETE CASCADE
);