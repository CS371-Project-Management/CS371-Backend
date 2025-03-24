CREATE TABLE IF NOT EXISTS ordering_histories
(
    id              CHAR(36) PRIMARY KEY,
    quiz_history_id CHAR(36) NOT NULL,
    answer          TEXT     NOT NULL,
    `order`           INT      NOT NULL,
    result          BOOLEAN  NOT NULL,
    FOREIGN KEY (quiz_history_id) REFERENCES quiz_histories (id) ON DELETE CASCADE
);