INSERT OR IGNORE
INTO department (id, 'department_name', 'type', 'order')
VALUES (1, 'system', 'system', -10);

INSERT OR IGNORE 
INTO tg_user ('department_id', 'order', 'chat_id', 'name', 'username', 'role', 'status', 'commit_src') 
VALUES (1, -10, 1, 'system', 'system', 'USER', 'system', 'manual');
