CREATE TABLE IF NOT EXISTS choice_quizzes
(
    quiz_id  CHAR(36)                    PRIMARY KEY ,
    question TEXT                        NOT NULL,
    type     ENUM ('single', 'multiple') NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id) ON DELETE CASCADE
);