-- "user" definition

CREATE TABLE "tg_user" (
	id INTEGER NOT NULL,
	department_id INTEGER NOT NULL,
	"order" INTEGER NOT NULL,
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
	"order" INTEGER NOT NULL, -- need to make it unique
	CONSTRAINT department_pk PRIMARY KEY (id)
);

-- "auth_user" definition

CREATE TABLE "auth_user" (
	id INTEGER NOT NULL,
	chat_id INTEGER NOT NULL,
	username TEXT NOT NULL,
	CONSTRAINT auth_key_pk PRIMARY KEY (id)
);

-- "invitation" definition

CREATE TABLE "invitation" (
	id INTEGER NOT NULL,
	code TEXT NOT NULL UNIQUE,
	created_by_user_id INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL,
	used_by_user_id INTEGER,
	used_at TIMESTAMP,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	CONSTRAINT invitation_pk PRIMARY KEY (id),
	CONSTRAINT invitation_created_by_fk FOREIGN KEY (created_by_user_id) REFERENCES auth_user(id),
	CONSTRAINT invitation_used_by_fk FOREIGN KEY (used_by_user_id) REFERENCES auth_user(id)
);
