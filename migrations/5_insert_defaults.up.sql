-- password: rootroot
INSERT INTO users (user_id, login, password) VALUES (0, 'root', '$2a$10$f9PsJKZKQeE.0JEz8qoRtuTUVQZB1nnSP7SeLx0KpYmj8w1JZb0sC');
INSERT INTO services (service_id, name, owner_fk) VALUES (0, 'ssonny', 0);
INSERT INTO signups (user_fk, service_fk) VALUES (0, 0);