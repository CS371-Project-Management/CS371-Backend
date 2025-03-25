CREATE TABLE IF NOT EXISTS quiz_histories
(
    id        CHAR(36) PRIMARY KEY                                                NOT NULL,
    quiz_id   CHAR(36)                                                            NOT NULL UNIQUE,
    user_id   CHAR(36)                                                            NOT NULL,
    quiz_type ENUM ('choice', 'ordering', 'missing_words')                        NOT NULL,
    note      TEXT,
    status    ENUM ('in_progress', 'in_correct', 'correct') DEFAULT 'in_progress' NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id) ON DELETE CASCADE
);