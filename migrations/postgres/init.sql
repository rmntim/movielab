CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL
);

CREATE TYPE sex AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS actors
(
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    sex       sex          NOT NULL,
    birthdate DATE         NOT NULL
);

CREATE TABLE IF NOT EXISTS movies
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(150) NOT NULL,
    description  VARCHAR(1000),
    release_date DATE         NOT NULL,
    rating       INT          NOT NULL
);

CREATE TABLE IF NOT EXISTS actors_movies
(
    id       SERIAL PRIMARY KEY,
    actor_id INT REFERENCES actors,
    movie_id INT REFERENCES movies,
    UNIQUE (actor_id, movie_id)
);
