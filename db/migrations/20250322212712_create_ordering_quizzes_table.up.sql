CREATE TABLE IF NOT EXISTS ordering_quizzes
(
    id       CHAR(36) PRIMARY KEY,
    quiz_id  CHAR(36) NOT NULL,
    question TEXT     NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id)
);