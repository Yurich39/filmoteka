package repo

import (
	"context"
	"database/sql"
	"fmt"

	"filmoteka/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ActorsMoviesRepo struct {
	db *sqlx.DB
}

func NewActorsMoviesRepo(db *sql.DB) *ActorsMoviesRepo {
	return &ActorsMoviesRepo{db: sqlx.NewDb(db, "postgres")}
}

const ActorMovieQuerySave = `INSERT INTO actors_movies (actor_id, movie_id) VALUES ($1, $2)`

func (r *ActorsMoviesRepo) Save(ctx context.Context, data entity.ActorMovie) error {

	//Проверяем наличие актера в таблице actors
	ActorQuery := `SELECT COUNT(*) FROM actors WHERE id = $1`

	var count1 int

	err := r.db.GetContext(ctx, &count1, ActorQuery,
		data.Actor_id,
	)

	if err != nil {
		return err
	}

	if count1 == 0 {
		return fmt.Errorf("actor_id was NOT found in database table 'actors'")
	}

	//Проверяем наличие фильма в таблице movies
	MovieQuery := `SELECT COUNT(*) FROM movies WHERE id = $1`

	var count2 int

	err = r.db.GetContext(ctx, &count2, MovieQuery,
		data.Movie_id,
	)

	if err != nil {
		return err
	}

	if count2 == 0 {
		return fmt.Errorf("movie_id was NOT found in database table 'movies'")
	}

	// Вносим данные в базу данных в таблицу movie_actors
	_, err = r.db.ExecContext(ctx, ActorMovieQuerySave,
		data.Actor_id,
		data.Movie_id,
	)

	if err != nil {
		return err
	}

	return nil
}

const ListActorsAndMoviesQuery = `SELECT 
			actors.id AS actor_id,
			actors.name AS actor_name,
			actors.surname AS actor_surname,
			movies.id AS movie_id,
			movies.title AS movie_title
			FROM actors
			JOIN 
			actors_movies ON actors.id = actors_movies.actor_id
			JOIN
			movies ON actors_movies.movie_id = movies.id`

func (r *ActorsMoviesRepo) List(ctx context.Context) ([]entity.ActorMovieData, error) {
	var data []entity.ActorMovieData
	err := r.db.SelectContext(ctx, &data, ListActorsAndMoviesQuery)

	if err != nil {
		return nil, fmt.Errorf("%s: DB returned error: %w", op, err)
	}

	return data, nil
}
