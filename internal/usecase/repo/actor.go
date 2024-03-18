package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"filmoteka/internal/controller/middleware/filter"
	"filmoteka/internal/controller/middleware/pagination"
	"filmoteka/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const op = "internal.usecase.repo"
const pageSize uint64 = 10

type ActorsRepo struct {
	db *sqlx.DB
}

func NewActorsRepo(db *sql.DB) *ActorsRepo {
	return &ActorsRepo{db: sqlx.NewDb(db, "postgres")}
}

const ActorQueryFind = `SELECT *, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth FROM actors WHERE  id = $1`

func (r *ActorsRepo) Get(ctx context.Context, id int) (entity.Actor, error) {

	var res entity.Actor
	err := r.db.GetContext(ctx, &res, ActorQueryFind, id)

	if err != nil {
		return entity.Actor{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

const ActorQuerySave = `INSERT INTO actors(name, surname, patronymic, gender, date_of_birth)
					VALUES($1, $2, $3, $4, TO_DATE($5, 'DD.MM.YYYY'))
					ON CONFLICT (name, surname) DO NOTHING
    				RETURNING id`

func (r *ActorsRepo) Save(ctx context.Context, data entity.ActorData) (int, error) {

	var res int

	err := r.db.GetContext(ctx, &res, ActorQuerySave,
		data.Name,
		data.Surname,
		data.Patronymic,
		data.Gender,
		data.DateOfBirth,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *ActorsRepo) Update(ctx context.Context, updates entity.Actor) (entity.Actor, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Update("actors")

	// Составим выражение для оператора SQL SET
	data, err := getMapActor(updates)

	if err != nil {
		return entity.Actor{}, fmt.Errorf("%s: Error: %w", op, err)
	}

	qb = qb.SetMap(data)

	// Составим выражение для оператора SQL Where
	stmt := fmt.Sprintf(" WHERE id = %d", *updates.Id)

	sql, i, err := qb.ToSql()

	sql = sql + stmt + " RETURNING id, name, surname, patronymic, gender, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth"

	if err != nil {
		return entity.Actor{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	var res entity.Actor

	err = r.db.GetContext(ctx, &res, sql, i...)

	if err != nil {
		return entity.Actor{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *ActorsRepo) Delete(ctx context.Context, id int) (entity.Actor, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Delete("actors")

	// Составим выражение для оператора SQL Where
	stmt := fmt.Sprintf(" WHERE id = %d", id)

	sql, i, err := qb.ToSql()

	sql = sql + stmt + " RETURNING id, name, surname, patronymic, gender, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth"

	if err != nil {
		return entity.Actor{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	var res entity.Actor

	err = r.db.GetContext(ctx, &res, sql, i...)

	if err != nil {
		return entity.Actor{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *ActorsRepo) List(ctx context.Context) ([]entity.Actor, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, patronymic, gender, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth").From("actors")

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

	// Используем оператор ORDER BY
	qb = qb.OrderBy("id ASC")

	qb = qb.Limit(pageSize)

	sql, i, err := qb.ToSql()

	if err != nil {
		return []entity.Actor{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	res := []entity.Actor{}
	err = r.db.SelectContext(ctx, &res, sql, i...)

	if err != nil {
		return []entity.Actor{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *ActorsRepo) Next(ctx context.Context) ([]entity.Actor, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, patronymic, gender, TO_CHAR(date_of_birth, 'DD.MM.YYYY') AS date_of_birth").From("actors")

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
		return []entity.Actor{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	res := []entity.Actor{}
	err = r.db.SelectContext(ctx, &res, sql, i...)

	if err != nil {
		return []entity.Actor{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func getMapActor(updates entity.Actor) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	if val := updates.ActorData.Name; val != nil {
		res["name"] = *val
	}

	if val := updates.ActorData.Surname; val != nil {
		res["surname"] = *val
	}

	if val := updates.ActorData.Patronymic; val != nil {
		res["patronymic"] = *val
	}

	if val := updates.ActorData.Gender; val != nil {
		res["gender"] = *val
	}

	if val := updates.ActorData.DateOfBirth; val != nil {
		layout := "02.01.2006" // шаблон формата даты в строке

		// распарсим строку в формат даты
		date, err := time.Parse(layout, *val)
		if err != nil {
			return res, fmt.Errorf("%s: Ошибка при парсинге даты", err)
		}

		// Преобразовать дату в формат PostgreSQL
		res["date_of_birth"] = date.Format("2006-01-02")
	}

	if len(res) == 0 {
		return res, fmt.Errorf("%s: Data for update operation were NOT specified", op)
	}

	return res, nil
}
