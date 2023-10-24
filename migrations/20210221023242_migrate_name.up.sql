-- CREATE TABLE IF NOT EXISTS people(
--     id serial PRIMARY KEY,
--     source VARCHAR(255),
--     destination VARCHAR(255),
--     original VARCHAR(255),
--     translation VARCHAR(255)
-- );

CREATE TABLE IF NOT EXISTS people(
	    id SERIAL PRIMARY KEY,
	    name TEXT NOT NULL,
	    surname TEXT NOT NULL,
		patronymic TEXT NOT NULL,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		nationality TEXT NOT NULL
);