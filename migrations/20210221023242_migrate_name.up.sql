CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS actors (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	surname VARCHAR(50) NOT NULL,
	patronymic VARCHAR(50),
    gender gender,
	date_of_birth DATE,
    CONSTRAINT unique_name_surname UNIQUE (name, surname)
);

CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description VARCHAR(1000),
    release_date DATE,
    rating INTEGER CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE IF NOT EXISTS actors_movies (
    movie_id INT REFERENCES movies(id),
    actor_id INT REFERENCES actors(id),
    PRIMARY KEY (movie_id, actor_id)
);
