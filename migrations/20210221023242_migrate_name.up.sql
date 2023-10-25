CREATE TABLE IF NOT EXISTS people(
	    id SERIAL PRIMARY KEY,
	    name TEXT NOT NULL,
	    surname TEXT NOT NULL,
		patronymic TEXT,
		age INTEGER,
		gender TEXT,
		nationality TEXT
);