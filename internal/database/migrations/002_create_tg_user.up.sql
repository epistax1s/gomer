PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS tg_user (
    id INTEGER NOT NULL,
    'department_id' INTEGER NOT NULL,
    'order' INTEGER NOT NULL,
    'chat_id' INTEGER NOT NULL,
    'redmine_id' INTEGER, 
    'name' TEXT NOT NULL,
    'username' TEXT NOT NULL,
    'role' TEXT NOT NULL,
    'status' TEXT NOT NULL,
    'commit_src' TEXT NOT NULL,
    'ee' INTEGER NOT NULL,
    CONSTRAINT chat_id_unique UNIQUE ('chat_id'),
    CONSTRAINT user_pk PRIMARY KEY (id),
    CONSTRAINT user_department_fk FOREIGN KEY ('department_id') REFERENCES department(id)
);
