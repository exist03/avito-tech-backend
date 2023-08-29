CREATE DATABASE avito;
GRANT ALL PRIVILEGES ON DATABASE avito TO postgres;
\c avito
DROP TABLE IF EXISTS accordance;
DROP TABLE IF EXISTS segment;
DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    id   INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);
CREATE TABLE segment
(
    id   INT PRIMARY KEY,
    name VARCHAR(50)
);
CREATE TABLE accordance
(
    user_id    INT REFERENCES users (id) ON DELETE CASCADE,
    segment_id INT REFERENCES segment (id) ON DELETE CASCADE,
    expires    BIGINT
);
CREATE TABLE history
(
    user_id    INT,
    segment_id INT,
    type       BOOL,
    time       BIGINT
);
INSERT INTO users(id, name) VALUES (1, 'name48'), (2, 'name2'), (3, 'name3'), (4, 'name4'), (5, 'name5'), (6, 'name6'), (7, 'name7'), (8, 'name8'), (9, 'name9'), (10, 'name10'), (11, 'name11'), (12, 'name12'), (13, 'name13'), (14, 'name14'), (15, 'name15'), (16, 'name16'), (17, 'name17'), (18, 'name18'), (19, 'name19');
