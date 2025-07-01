PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS invitation (
    id INTEGER NOT NULL,
    'code' TEXT NOT NULL UNIQUE,
    'created_by_user_id' INTEGER NOT NULL,
    'created_at' TEXT NOT NULL,
    'used_by_user_id' INTEGER,
    'used_at' TEXT,
    'is_active' INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT invitation_pk PRIMARY KEY (id),
    CONSTRAINT invitation_created_by_fk FOREIGN KEY ('created_by_user_id') REFERENCES tg_user(id),
    CONSTRAINT invitation_used_by_fk FOREIGN KEY ('used_by_user_id') REFERENCES tg_user(id)
);
