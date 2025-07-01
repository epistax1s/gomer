PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS 'commit' (
    id INTEGER NOT NULL,
    'user_id' INTEGER,
    'commit_payload' TEXT NOT NULL,
    'commit_date' TEXT NOT NULL,
    CONSTRAINT commit_pk PRIMARY KEY (id),
    CONSTRAINT commit_tg_user_FK FOREIGN KEY ('user_id') REFERENCES tg_user(id)
);
