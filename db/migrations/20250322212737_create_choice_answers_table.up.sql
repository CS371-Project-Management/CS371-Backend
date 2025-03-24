CREATE TABLE IF NOT EXISTS choice_answers
(
    id      CHAR(36) PRIMARY KEY,
    quiz_id CHAR(36) NOT NULL,
    answer  TEXT     NOT NULL,
    result  BOOLEAN  NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES choice_quizzes (quiz_id) ON DELETE CASCADE
);