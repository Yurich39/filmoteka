package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"filmoteka/internal/controller/middleware/filter"
	"filmoteka/internal/controller/middleware/pagination"
	"filmoteka/internal/controller/middleware/sort"
	"filmoteka/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type MoviesRepo struct {
	db *sqlx.DB
}

func NewMoviesRepo(db *sql.DB) *MoviesRepo {
	return &MoviesRepo{db: sqlx.NewDb(db, "postgres")}
}

const MovieQueryFind = `SELECT *, TO_CHAR(release_date, 'DD.MM.YYYY') AS release_date FROM movies WHERE id = $1`

func (r *MoviesRepo) Get(ctx context.Context, id int) (entity.Movie, error) {

	var res entity.Movie
	err := r.db.GetContext(ctx, &res, MovieQueryFind, id)

	if err != nil {
		return entity.Movie{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

const MovieQueryFindMovie = `
SELECT DISTINCT movies.id, movies.title, movies.description, TO_CHAR(movies.release_date, 'DD.MM.YYYY') AS release_date, movies.rating
FROM movies
JOIN actors_movies ON movies.id = actors_movies.movie_id
JOIN actors ON actors_movies.actor_id = actors.id
WHERE ($1 = '' OR LOWER(movies.title) LIKE '%' || LOWER($1) || '%')
AND ($2 = '' OR LOWER(actors.name) LIKE '%' || LOWER($2) || '%')`

func (r *MoviesRepo) GetMovie(ctx context.Context) ([]entity.Movie, error) {

	filter_options, _ := ctx.Value(filter.FilterOptionsContextKey).(map[string][]string)
	titles := filter_options["title"]
	actors := filter_options["actor_name"]

	var title, actor string
	if len(titles) == 0 {
		title = ""
	} else {
		title = titles[0]
	}

	if len(actors) == 0 {
		actor = ""
	} else {
		actor = actors[0]
	}

	var res []entity.Movie
	err := r.db.SelectContext(ctx, &res, MovieQueryFindMovie, title, actor)

	if err != nil {
		return []entity.Movie{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

const MovieQuerySave = `INSERT INTO movies(title, description, release_date, rating)
					VALUES($1, $2, TO_DATE($3, 'DD.MM.YYYY'), $4)
					RETURNING id`

func (r *MoviesRepo) Save(ctx context.Context, data entity.MovieData) (int, error) {

	var res int

	err := r.db.GetContext(ctx, &res, MovieQuerySave,
		data.Title,
		data.Description,
		data.ReleaseDate,
		data.Rating,
	)

	if err != nil {
		return 0, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *MoviesRepo) Update(ctx context.Context, updates entity.Movie) (entity.Movie, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Update("movies")

	// Составим выражение для оператора SQL SET
	data, err := getMapMovie(updates)

	if err != nil {
		return entity.Movie{}, fmt.Errorf("%s: Error: %w", op, err)
	}

	qb = qb.SetMap(data)

	// Составим выражение для оператора SQL Where
	stmt := fmt.Sprintf(" WHERE id = %d", *updates.Id)

	sql, i, err := qb.ToSql()

	sql = sql + stmt + " RETURNING id, title, description, TO_CHAR(release_date, 'DD.MM.YYYY') AS release_date, rating"

	if err != nil {
		return entity.Movie{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	var res entity.Movie

	err = r.db.GetContext(ctx, &res, sql, i...)

	if err != nil {
		return entity.Movie{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *MoviesRepo) Delete(ctx context.Context, id int) (entity.Movie, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Delete("movies")

	// Составим выражение для оператора SQL Where
	stmt := fmt.Sprintf(" WHERE id = %d", id)

	sql, i, err := qb.ToSql()

	sql = sql + stmt + " RETURNING id, title, description, TO_CHAR(release_date, 'DD.MM.YYYY') AS release_date, rating"

	if err != nil {
		return entity.Movie{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	var res entity.Movie

	err = r.db.GetContext(ctx, &res, sql, i...)

	if err != nil {
		return entity.Movie{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *MoviesRepo) List(ctx context.Context) ([]entity.Movie, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, title, description, rating, TO_CHAR(release_date, 'DD.MM.YYYY') AS release_date").From("movies")

	// Составим выражение для оператора SQL Where ... AND ...
	filter_options, _ := ctx.Value(filter.FilterOptionsContextKey).(map[string][]string)

	stmt := []string{}
	for k, v := range filter_options {
		for _, val := range v {
			stmt = append(stmt, fmt.Sprintf("%s = '%s'", k, val))
		}
	}

	// Решаем, использовать оператор WHERE или нет
	if len(stmt) != 0 {
		qb = qb.Where(strings.Join(stmt, " AND "))
	}

	// Используем оператор ORDER BY для сортировки
	sort_options, _ := ctx.Value(sort.SortOptionsContextKey).(map[string]string)

	stmt = []string{}

	for k, v := range sort_options {
		stmt = append(stmt, fmt.Sprintf("%s %s", k, v))
	}

	str := strings.Join(stmt, ", ")

	if len(str) == 0 {
		qb = qb.OrderBy("rating DESC")
	} else {
		qb = qb.OrderBy(str)
	}

	qb = qb.Limit(pageSize)

	sql, i, err := qb.ToSql()

	if err != nil {
		return []entity.Movie{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	res := []entity.Movie{}
	err = r.db.SelectContext(ctx, &res, sql, i...)

	if err != nil {
		return []entity.Movie{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *MoviesRepo) Next(ctx context.Context) ([]entity.Movie, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, patronymic, gender, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth").From("movies")

	// Составим выражение для оператора SQL Where ... AND ...
	filter_options, _ := ctx.Value(filter.FilterOptionsContextKey).(map[string][]string)

	stmt := []string{}
	for k, v := range filter_options {
		for _, val := range v {
			stmt = append(stmt, fmt.Sprintf("%s = '%s'", k, val))
		}
	}

	// Решаем, использовать оператор WHERE или нет
	if len(stmt) != 0 {
		qb = qb.Where(strings.Join(stmt, " AND "))
	}

	// Условие пагинации
	personID := ctx.Value(pagination.NextPersonID).(int)
	qb = qb.Where(fmt.Sprintf("id >= %d", personID))

	// Используем оператор ORDER BY
	qb = qb.OrderBy("id ASC")

	qb = qb.Limit(pageSize)

	sql, i, err := qb.ToSql()

	if err != nil {
		return []entity.Movie{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	res := []entity.Movie{}
	err = r.db.SelectContext(ctx, &res, sql, i...)

	if err != nil {
		return []entity.Movie{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func getMapMovie(updates entity.Movie) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	if val := updates.MovieData.Title; val != nil {
		res["title"] = *val
	}

	if val := updates.MovieData.Description; val != nil {
		res["description"] = *val
	}

	if val := updates.MovieData.ReleaseDate; val != nil {
		layout := "02.01.2006" // шаблон формата даты в строке

		// распарсим строку в формат даты
		date, err := time.Parse(layout, *val)
		if err != nil {
			return res, fmt.Errorf("%s: Ошибка при парсинге даты", err)
		}

		// Преобразовать дату в формат PostgreSQL
		res["release_date"] = date.Format("2006-01-02")
	}

	if val := updates.Rating; val != nil {
		res["rating"] = *val
	}

	if len(res) == 0 {
		return res, fmt.Errorf("%s: Data for update operation were NOT specified", op)
	}

	return res, nil
}
