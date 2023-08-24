CREATE DATABASE avito;
GRANT ALL PRIVILEGES ON DATABASE avito TO postgres;
\c avito
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS segment;
DROP TABLE IF EXISTS accordance;

CREATE TABLE users
(
    id   INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);
CREATE TABLE segment
(
    id   INT PRIMARY KEY ,
    name VARCHAR(50)
);
CREATE TABLE accordance
(
    user_id    INT REFERENCES users (id),
    segment_id INT REFERENCES segment (id)
);