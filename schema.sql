-- "user" definition

CREATE TABLE "tg_user" (
	id INTEGER NOT NULL,
	department_id INTEGER NOT NULL,
	order INTEGER NOT NULL, -- need to make it unique
	chat_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	username TEXT NOT NULL,
	role TEXT NOT NULL,
	status TEXT NOT NULL,
	CONSTRAINT user_pk PRIMARY KEY (id),
	CONSTRAINT user_department_fk FOREIGN KEY (department_id) REFERENCES department(id)
);

-- "commit" definition

CREATE TABLE "commit" (
	id INTEGER NOT NULL,
	user_id INTEGER,
	commit_payload TEXT NOT NULL,
	commit_date TEXT NOT NULL,
	CONSTRAINT commit_pk PRIMARY KEY(id)
	CONSTRAINT commit_tg_user_FK FOREIGN KEY (user_id) REFERENCES tg_user(id)
);

-- "tg_group" definition

CREATE TABLE "tg_group" (
	id INTEGER NOT NULL,
	group_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	CONSTRAINT tg_group_pk PRIMARY KEY (id)
);


-- "department" definition

CREATE TABLE "department" (
	id INTEGER NOT NULL,
	department_name TEXT NOT NULL,
	order INTEGER NOT NULL, -- need to make it unique
	CONSTRAINT department_pk PRIMARY KEY (id)
);