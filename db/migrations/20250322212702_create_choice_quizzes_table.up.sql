CREATE TABLE IF NOT EXISTS choice_quizzes
(
    id       CHAR(36) PRIMARY KEY,
    quiz_id  CHAR(36)                    NOT NULL,
    question TEXT                        NOT NULL,
    type     ENUM ('single', 'multiple') NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id)
);