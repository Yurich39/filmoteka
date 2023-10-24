package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"people-finder/internal/entity"

	"github.com/Masterminds/squirrel"
)

const op = "internal.usecase.repo"

type PersonRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *PersonRepo {
	return &PersonRepo{db: db}
}

const QuerySave = `INSERT INTO people(name, surname, patronymic, age, gender, nationality)
					VALUES($1, $2, $3, $4, $5, $6)
					RETURNING id, name, surname, patronymic, age, gender, nationality`

func (r *PersonRepo) Save(ctx context.Context, person entity.Person) ([]entity.Person, error) {

	rows, err := r.db.Query(QuerySave,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	res, err := formResult(rows)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: Error during processing of DB result: %w", op, err)
	}

	return res, nil
}

func (r *PersonRepo) Update(ctx context.Context, updates entity.UpdateRequest) ([]entity.Person, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Update("people")

	// Составим выражение для оператора SQL SET
	qb = qb.SetMap(updates.NewFields)

	// Составим выражение для оператора SQL Where ... AND ...
	stmt := []string{}
	for k, v := range updates.Filters {
		stmt = append(stmt, fmt.Sprintf("%s = '%s'", k, v))
	}

	// Решаем, использовать оператор WHERE или нет
	if len(stmt) != 0 {
		qb = qb.Where(strings.Join(stmt, " AND "))
	}

	sql, i, err := qb.ToSql()

	sql = sql + "RETURNING id, name, surname, patronymic, age, gender, nationality"

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	rows, err := r.db.Query(sql, i...)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	res, err := formResult(rows)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: Error during processing of DB result: %w", op, err)
	}

	return res, nil
}

func (r *PersonRepo) Delete(ctx context.Context, deleter entity.DelRequest) ([]entity.Person, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Delete("people")

	// Составим выражение для оператора SQL Where ... AND ...
	stmt := []string{}
	for k, v := range deleter.Filters {
		stmt = append(stmt, fmt.Sprintf("%s = '%s'", k, v))
	}

	// Решаем, использовать оператор WHERE или нет
	if len(stmt) != 0 {
		qb = qb.Where(strings.Join(stmt, " AND "))
	}

	sql, i, err := qb.ToSql()

	sql = sql + "RETURNING id, name, surname, patronymic, age, gender, nationality"

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: squirrel failed to build sql statement : %w", op, err)
	}

	rows, err := r.db.Query(sql, i...)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	res, err := formResult(rows)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: Error during processing of DB result: %w", op, err)
	}

	return res, nil
}

const pageSize uint64 = 3

func (r *PersonRepo) Find(ctx context.Context, getter entity.Options) ([]entity.Person, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, patronymic, age, gender, nationality").From("people")

	// Составим выражение для оператора SQL Where ... AND ...
	And := make(map[string]interface{})
	for k, v := range getter.Where {
		And[k] = v[0]
	}

	stmt := []string{}
	for k, v := range And {
		stmt = append(stmt, fmt.Sprintf("%s = '%s'", k, v))
	}

	// Решаем, использовать оператор WHERE или нет
	if len(stmt) != 0 {
		qb = qb.Where(strings.Join(stmt, " AND "))
	}

	// Составим выражение для оператора SQL ORDER BY
	stmt = []string{}
	for i := range getter.OrderBy {
		stmt = append(stmt, fmt.Sprintf("%s %s", getter.OrderBy[i], getter.Order[i]))
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

	rows, err := r.db.Query(sql, i...)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	res, err := formResult(rows)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: Error during processing of DB result: %w", op, err)
	}

	return res, nil
}

func (r *PersonRepo) ListPeople(ctx context.Context, getter entity.Options, personID int) ([]entity.Person, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Подготавливаем SQL запрос
	qb := psql.Select("id, name, surname, patronymic, age, gender, nationality").From("people")

	// Составим выражение для оператора SQL Where ... AND ...
	And := make(map[string]interface{})
	for k, v := range getter.Where {
		And[k] = v[0]
	}

	stmt := []string{}
	for k, v := range And {
		stmt = append(stmt, fmt.Sprintf("%s = '%s'", k, v))
	}

	// Решаем, использовать оператор WHERE или нет
	if len(stmt) != 0 {
		qb = qb.Where(strings.Join(stmt, " AND "))
	}

	qb = qb.Where(fmt.Sprintf("id >= %d", personID))

	// Составим выражение для оператора SQL ORDER BY
	stmt = []string{}
	for i := range getter.OrderBy {
		stmt = append(stmt, fmt.Sprintf("%s %s", getter.OrderBy[i], getter.Order[i]))
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

	rows, err := r.db.Query(sql, i...)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: DB method 'Query' returned error: %w", op, err)
	}

	res, err := formResult(rows)

	if err != nil {
		return []entity.Person{}, fmt.Errorf("%s: Error during processing of DB result: %w", op, err)
	}

	return res, nil
}

func formResult(rows *sql.Rows) ([]entity.Person, error) {
	res := []entity.Person{}

	for rows.Next() {
		var DBResponse entity.Person

		err := rows.Scan(&DBResponse.ID, &DBResponse.Name, &DBResponse.Surname, &DBResponse.Patronymic, &DBResponse.Age, &DBResponse.Gender, &DBResponse.Nationality)

		if err != nil {
			return []entity.Person{}, fmt.Errorf("%s: failed to unmarshal DB response: %w", op, err)
		}

		res = append(res, DBResponse)
	}

	if len(res) == 0 {
		return []entity.Person{}, fmt.Errorf("%s: Data not found - result is empty", op)
	}

	return res, nil
}

// func (r *PersonRepo) Update(ctx context.Context, person entity.Person) ([]entity.Person, error) {
// 	if _, err := r.db.ExecContext(ctx, QuerySave,
// 		person.Age,
// 		person.Gender,
// 		person.Nationality,
// 		person.Name,
// 		person.Surname,
// 		person.Patronymic,
// 	); err != nil {
// 		return fmt.Errorf("%s: Method 'Save' returned error: %w", op, err)
// 	}

// 	return nil
// }
