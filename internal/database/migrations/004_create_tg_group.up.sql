PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS tg_group (
    id INTEGER NOT NULL,
    'group_id' INTEGER NOT NULL,
    'title' TEXT NOT NULL,
    CONSTRAINT tg_group_pk PRIMARY KEY (id)
);
