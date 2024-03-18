CREATE TYPE role AS ENUM ('user', 'admin');

CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL,
    role     role                NOT NULL
);

CREATE TYPE sex AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS movies
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(150) NOT NULL,
    description  VARCHAR(1000),
    release_date TIMESTAMPTZ  NOT NULL,
    rating       INT          NOT NULL
        CONSTRAINT rating_check CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE IF NOT EXISTS actors
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    sex        sex          NOT NULL,
    birth_date TIMESTAMPTZ  NOT NULL
);

CREATE TABLE IF NOT EXISTS movie_actors
(
    movie_id INT NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
    actor_id INT NOT NULL REFERENCES actors (id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, actor_id)
);