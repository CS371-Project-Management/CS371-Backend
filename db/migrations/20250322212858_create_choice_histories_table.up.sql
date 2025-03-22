CREATE TABLE IF NOT EXISTS choice_histories
(
    id              CHAR(36) PRIMARY KEY,
    quiz_history_id CHAR(36) NOT NULL,
    answer          TEXT     NOT NULL,
    result          BOOLEAN  NOT NULL,
    FOREIGN KEY (quiz_history_id) REFERENCES quiz_histories (id)
);