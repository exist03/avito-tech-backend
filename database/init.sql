CREATE
DATABASE avito;
GRANT ALL PRIVILEGES ON DATABASE
avito TO postgres;
\c
avito
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
