DROP DATABASE IF EXISTS goblog;
CREATE DATABASE goblog;

CREATE TABLE users (
    userid SERIAL PRIMARY KEY,
    name TEXT,
    age INT,
    location TEXT
);