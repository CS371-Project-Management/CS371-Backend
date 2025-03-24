CREATE TABLE IF NOT EXISTS ordering_quizzes
(
    quiz_id  CHAR(36) NOT NULL,
    question TEXT     NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id)
);