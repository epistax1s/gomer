PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS department (
    id INTEGER NOT NULL,
    'department_name' TEXT NOT NULL,
    'type' TEXT NOT NULL, 
    'order' INTEGER NOT NULL,
    CONSTRAINT department_pk PRIMARY KEY (id),
    CONSTRAINT department_order_unique UNIQUE ('order')
);