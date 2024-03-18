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
);

CREATE TABLE IF NOT EXISTS actors
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    sex        sex          NOT NULL,
    birth_date TIMESTAMPTZ  NOT NULL,
    movie_id   INT          NOT NULL REFERENCES movies (id)
);

CREATE VIEW movies_actors AS
SELECT movies.*,
       (SELECT array_to_json(array_agg(row_to_json(actorslist.*))) as array_to_json
        FROM (SELECT id, name, sex, birth_date
              FROM actors
              where movie_id = movies.id) actorslist) as actors
FROM movies;