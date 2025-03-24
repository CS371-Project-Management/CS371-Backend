CREATE TABLE IF NOT EXISTS missing_words_quizzes
(
    quiz_id  CHAR(36) NOT NULL,
    question TEXT     NOT NULL,
    answer   TEXT     NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id)
);