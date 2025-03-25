CREATE TABLE IF NOT EXISTS missing_words_quizzes
(
    quiz_id  CHAR(36) PRIMARY KEY,
    question TEXT     NOT NULL,
    answer   TEXT     NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id) ON DELETE CASCADE
);