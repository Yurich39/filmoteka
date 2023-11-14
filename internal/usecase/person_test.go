package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"people-finder/internal/entity"
	"people-finder/internal/usecase"
)

var errInternalServErr = errors.New("internal server error")

func person(t *testing.T) (*usecase.PersonUseCase, *MockPersonRepo, *MockEnrichWebAPI) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockPersonRepo(mockCtl)
	webAPI := NewMockEnrichWebAPI(mockCtl)
	log := NewMockLogger(mockCtl)

	testPersonUseCase := usecase.New(repo, webAPI, log)

	return testPersonUseCase, repo, webAPI
}

type testFind struct {
	name string
	val  int
	mock func()
	res  interface{}
	err  error
}

func TestFind(t *testing.T) {
	t.Parallel()

	person, repo, _ := person(t)

	tests := []testFind{
		{
			name: "id is out of range",
			val:  1000,
			mock: func() {
				repo.EXPECT().Get(context.Background(), 1000).Return(entity.Person{}, errInternalServErr)
			},
			res: entity.Person{},
			err: errInternalServErr,
		},
		{
			name: "existing id",
			val:  1,
			mock: func() {
				repo.EXPECT().Get(context.Background(), 1).Return(entity.Person{}, nil)
			},
			res: entity.Person{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			res, err := person.Find(context.Background(), tc.val)

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

type testSave struct {
	name  string
	input func() entity.Data
	mock  func(data entity.Data)
	res   func(data entity.Data) entity.Person
	err   error
}

func TestSave(t *testing.T) {
	t.Parallel()

	person, repo, webAPI := person(t)

	// Тестовые данные
	tests := []testSave{
		{
			name: "ok result",
			input: func() entity.Data {

				// Создаем структуру entity.Data из строки
				s := `{
						"name": "Leo",
						"surname": "Kim"
					}`

				var data entity.Data
				r := strings.NewReader(s)
				render.DecodeJSON(r, &data)

				return data
			},
			mock: func(data entity.Data) {
				webAPI.EXPECT().EnrichAge(*data.Name).Return(1, nil)
				webAPI.EXPECT().EnrichGender(*data.Name).Return("", nil)
				webAPI.EXPECT().EnrichNationality(*data.Name).Return("", nil)

				// Enrich data
				age := 1
				gender := ""
				nationality := ""

				data.Age = &age
				data.Gender = &gender
				data.Nationality = &nationality
				repo.EXPECT().Save(context.Background(), data).Return(1, nil)
			},
			res: func(data entity.Data) entity.Person {
				// Создаем структуру entity.Person
				id := 1
				age := 1
				gender := ""
				nationality := ""

				person := entity.Person{
					Id:   &id,
					Data: data,
				}

				person.Age = &age
				person.Gender = &gender
				person.Nationality = &nationality

				return person
			},
			err: nil,
		},
		{
			name: "db error",
			input: func() entity.Data {

				// Создаем структуру entity.Data из строки
				s := `{
						"name": "Leo",
						"surname": "Kim"
					}`

				var data entity.Data
				r := strings.NewReader(s)
				render.DecodeJSON(r, &data)

				return data
			},
			mock: func(data entity.Data) {
				webAPI.EXPECT().EnrichAge(*data.Name).Return(1, nil)
				webAPI.EXPECT().EnrichGender(*data.Name).Return("", nil)
				webAPI.EXPECT().EnrichNationality(*data.Name).Return("", nil)

				// Enrich data
				age := 1
				gender := ""
				nationality := ""

				data.Age = &age
				data.Gender = &gender
				data.Nationality = &nationality
				repo.EXPECT().Save(context.Background(), data).Return(0, errInternalServErr)
			},
			res: func(data entity.Data) entity.Person {
				// Создаем структуру entity.Person
				return entity.Person{}
			},
			err: errInternalServErr,
		},
	}

	// Запуск тестов
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data := tc.input()

			tc.mock(data)

			res, err := person.Save(context.Background(), data)

			require.Equal(t, res, tc.res(data))
			require.ErrorIs(t, err, tc.err)
		})
	}
}

type testUpdate struct {
	name  string
	input func() entity.Person
	mock  func(person entity.Person)
	res   func(person entity.Person) entity.Person
	err   error
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	person, repo, _ := person(t)

	// Тестовые данные
	tests := []testUpdate{
		{
			name: "ok result",
			input: func() entity.Person {

				// Создаем структуру entity.Person из строки
				s := `{
						"id": 1,
						"patronymic": "Vitalevich",
						"nationality": "RU"
					}`

				var person entity.Person
				r := strings.NewReader(s)
				render.DecodeJSON(r, &person)

				return person
			},
			mock: func(person entity.Person) {
				// Создадим new person
				newPerson := person

				name := "Yuriy"
				surname := "Smith"
				age := 1
				gender := ""

				newPerson.Name = &name
				newPerson.Surname = &surname
				newPerson.Age = &age
				newPerson.Gender = &gender

				repo.EXPECT().Update(context.Background(), person).Return(newPerson, nil)
			},
			res: func(person entity.Person) entity.Person {
				// Создадим new person
				newPerson := person

				name := "Yuriy"
				surname := "Smith"
				age := 1
				gender := ""

				newPerson.Name = &name
				newPerson.Surname = &surname
				newPerson.Age = &age
				newPerson.Gender = &gender

				return newPerson
			},
			err: nil,
		},
		{
			name: "db error",
			input: func() entity.Person {

				// Создаем структуру entity.Person из строки
				s := `{
						"id": 1,
						"patronymic": "Vitalevich",
						"nationality": "RU"
					}`

				var person entity.Person
				r := strings.NewReader(s)
				render.DecodeJSON(r, &person)

				return person
			},
			mock: func(person entity.Person) {
				repo.EXPECT().Update(context.Background(), person).Return(entity.Person{}, errInternalServErr)
			},
			res: func(person entity.Person) entity.Person {
				return entity.Person{}
			},
			err: errInternalServErr,
		},
	}

	// Запуск тестов
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data := tc.input()

			tc.mock(data)

			res, err := person.Update(context.Background(), data)

			require.Equal(t, res, tc.res(data))
			require.ErrorIs(t, err, tc.err)
		})
	}
}

type testDelete struct {
	name string
	val  int
	mock func()
	res  interface{}
	err  error
}

func TestDelete(t *testing.T) {
	t.Parallel()

	person, repo, _ := person(t)

	tests := []testDelete{
		{
			name: "id is out of range",
			val:  1000,
			mock: func() {
				repo.EXPECT().Delete(context.Background(), 1000).Return(entity.Person{}, errInternalServErr)
			},
			res: entity.Person{},
			err: errInternalServErr,
		},
		{
			name: "existing id",
			val:  1,
			mock: func() {
				repo.EXPECT().Delete(context.Background(), 1).Return(entity.Person{}, nil)
			},
			res: entity.Person{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			res, err := person.Delete(context.Background(), tc.val)

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

type testList struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func TestList(t *testing.T) {
	t.Parallel()

	person, repo, _ := person(t)

	tests := []testList{
		{
			name: "db error",
			mock: func() {
				repo.EXPECT().List(context.Background()).Return([]entity.Person{}, errInternalServErr)
			},
			res: []entity.Person{},
			err: errInternalServErr,
		},
		{
			name: "ok",
			mock: func() {
				repo.EXPECT().List(context.Background()).Return([]entity.Person{}, nil)
			},
			res: []entity.Person{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			res, err := person.List(context.Background())

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

type testNext struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func TestNext(t *testing.T) {
	t.Parallel()

	person, repo, _ := person(t)

	tests := []testNext{
		{
			name: "db error",
			mock: func() {
				repo.EXPECT().Next(context.Background()).Return([]entity.Person{}, errInternalServErr)
			},
			res: []entity.Person{},
			err: errInternalServErr,
		},
		{
			name: "ok",
			mock: func() {
				repo.EXPECT().Next(context.Background()).Return([]entity.Person{}, nil)
			},
			res: []entity.Person{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			res, err := person.Next(context.Background())

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
