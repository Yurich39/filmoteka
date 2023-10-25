package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"people-finder/internal/controller/middleware/filter"
	"people-finder/internal/controller/middleware/pagination"
	"people-finder/internal/controller/middleware/sort"

	"people-finder/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const op = "internal.usecase.repo"

type PersonRepo struct {
	db *sqlx.DB
}

func New(db *sql.DB) *PersonRepo {
	return &PersonRepo{db: sqlx.NewDb(db, "postgres")}
}

const QueryFind = `SELECT * FROM people WHERE id = $1`

func (r *PersonRepo) Find(ctx context.Context, id entity.Id) ([]entity.Person, error) {

	var row entity.Person
	var res []entity.Person = []entity.Person{}

	err := r.db.GetContext(ctx, &row, QueryFind, id.Id)
	res = append(res, row)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

const QuerySave = `INSERT INTO people(name, surname, patronymic, age, gender, nationality)
					VALUES($1, $2, $3, $4, $5, $6)
					RETURNING id, name, surname, patronymic, age, gender, nationality`

func (r *PersonRepo) Save(ctx context.Context, person entity.Person) ([]entity.Person, error) {

	var row entity.Person
	var res []entity.Person = []entity.Person{}

	err := r.db.GetContext(ctx, &row, QuerySave,
		person.Data.Name,
		person.Data.Surname,
		person.Data.Patronymic,
		person.Data.Age,
		person.Data.Gender,
		person.Data.Nationality,
	)
	res = append(res, row)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *PersonRepo) Update(ctx context.Context, updates entity.Person) ([]entity.Person, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Update("people")

	// Составим выражение для оператора SQL SET
	data, err := getMap(updates)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: Error: %w", op, err)
	}

	qb = qb.SetMap(data)

	// Составим выражение для оператора SQL Where
	stmt := fmt.Sprintf(" WHERE id = %d", updates.Id.Id)

	sql, i, err := qb.ToSql()

	sql = sql + stmt + " RETURNING id, name, surname, patronymic, age, gender, nationality"

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	var row entity.Person
	var res []entity.Person = []entity.Person{}

	err = r.db.GetContext(ctx, &row, sql, i...)
	res = append(res, row)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func getMap(updates entity.Person) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	if val := updates.Data.Name; val != nil {
		res["name"] = *val
	}

	if val := updates.Data.Surname; val != nil {
		res["surname"] = *val
	}

	if val := updates.Data.Patronymic; val != nil {
		res["patronymic"] = *val
	}

	if val := updates.Data.Age; val != nil {
		res["age"] = *val
	}

	if val := updates.Data.Gender; val != nil {
		res["gender"] = *val
	}

	if val := updates.Data.Nationality; val != nil {
		res["nationality"] = *val
	}

	if len(res) == 0 {
		return res, fmt.Errorf("%s: Data for update operation were NOT specified", op)
	}

	return res, nil
}

func (r *PersonRepo) Delete(ctx context.Context, id entity.Id) ([]entity.Person, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Delete("people")

	// Составим выражение для оператора SQL Where
	stmt := fmt.Sprintf(" WHERE id = %d", id.Id)

	sql, i, err := qb.ToSql()

	sql = sql + stmt + " RETURNING id, name, surname, patronymic, age, gender, nationality"

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	var row entity.Person
	var res []entity.Person = []entity.Person{}

	err = r.db.GetContext(ctx, &row, sql, i...)
	res = append(res, row)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

const pageSize uint64 = 3

func (r *PersonRepo) List(ctx context.Context) ([]entity.Person, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, patronymic, age, gender, nationality").From("people")

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

	// Составим выражение для оператора SQL ORDER BY
	sort_by, _ := ctx.Value(sort.SortByContextKey).([]string)
	sort_order, _ := ctx.Value(sort.SortOrderContextKey).([]string)

	stmt = []string{}
	for i := 0; i < len(sort_by); i++ {
		stmt = append(stmt, fmt.Sprintf("%s %s", sort_by[i], sort_order[i]))
	}

	// Решаем, использовать оператор ORDER BY или нет
	if len(stmt) != 0 {
		qb = qb.OrderBy(strings.Join(stmt, ", "))
	}

	qb = qb.Limit(pageSize)

	sql, i, err := qb.ToSql()

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	res := []entity.Person{}
	err = r.db.SelectContext(ctx, &res, sql, i...)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}

func (r *PersonRepo) Next(ctx context.Context) ([]entity.Person, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, patronymic, age, gender, nationality").From("people")

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

	// Составим выражение для оператора SQL ORDER BY
	sort_by, _ := ctx.Value(sort.SortByContextKey).([]string)
	sort_order, _ := ctx.Value(sort.SortOrderContextKey).([]string)

	stmt = []string{}
	for i := 0; i < len(sort_by); i++ {
		stmt = append(stmt, fmt.Sprintf("%s %s", sort_by[i], sort_order[i]))
	}

	// Решаем, использовать оператор ORDER BY или нет
	if len(stmt) != 0 {
		qb = qb.OrderBy(strings.Join(stmt, ", "))
	}

	qb = qb.Limit(pageSize)

	sql, i, err := qb.ToSql()

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	res := []entity.Person{}
	err = r.db.SelectContext(ctx, &res, sql, i...)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	return res, nil
}
