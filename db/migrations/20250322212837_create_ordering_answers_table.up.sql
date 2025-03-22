CREATE TABLE IF NOT EXISTS ordering_answers
(
    id      CHAR(36) PRIMARY KEY,
    quiz_id CHAR(36) NOT NULL,
    answer  TEXT     NOT NULL,
    `order`   INT      NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES ordering_quizzes (id)
);