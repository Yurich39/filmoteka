# filmoteka

Проект написан по принципу "чистой" архитектуры.

## Описание

Проект filmoteka позволяет пользователям хранить информацию о фильмах и актерах, а так же получать эту информацию по запросу.
В качестве хранилища данных используется БД PostgreSQL.

Проект запускается в докер контейнерах: контейнер app с http-сервером и контейнер postgres с базой данных.

Реализованный в проекте web-сервер позволяет пользователю без аутентификации:
- 1.Осуществлять поиск актеров в БД по id актера
- 2.Получать из БД список актеров и список фильмов с их участием
- 3.Получать информацию о фильмах из БД как по id, так и по фрагменту названия фильма или фрагменту имени актера
- 4.Получать список фильмов из БД с возможностью сортировки по названию, рейтингу, дате выпуска

Реализованный в проекте web-сервер позволяет администратору (требуется регистрация):
- 1.Добавлять информацию в БД об актерах (имя, пол, дата рождения)
- 2.Изменять информацию об актерах (любое из полей или несколько полей)
- 3.Удалять информацию об актере по его id в БД
- 4.Добавлять фильмы и информацию о фильмах в БД (название фильма от 1 до 150 символов, описание не более 1000 символов, дата выпуска, рейтинг от 0 до 10)
- 5.Изменять информацию о фильмах (любое из полей или несколько полей)
- 6.А так же функции, доступные пользователям без аутентификации

__Алгоритм установки и запуска проекта:__
Проект упакован в два докер контейнера:
- **app** - web-сервер, написанный на Go 1.20
- **postgres** - postgresql db

Для запуска проекта используется make файл.

Запуск проекта в Docker контейнерах:
```sh
$ make compose-up
```
Так же в Makefile проекта доступны и прочие вспомагательные функции.

## База данных
База данных PostgreSQL.
Сервер работает только с одной БД.
Сервер работает с тремя таблицами: "actors", "movies", "actors_movies" в БД.

При первом запуске проекта путем миграции будет создана структура БД.
Между таблицами "actors" и "movies" установлена связь many-to-many с использованием вспомагательной таблицы "actors_movies".

Таблица "actors" состоит из следующих полей:
"id" (Pk) int, "name" text, "surname" text, "patronymic" text, "gender" text, "date_of_birth" date

Таблица "movies" состоит из следующих полей:
"id" (Pk) int, "title" text, "description" text, "release_date" date, "rating" int

Таблица "actors_movies" состоит из следующих полей:
"movie_id" (Fk) int, "actor_id" (FK) int

## `Логирование`
Логирование с использованием slog. Уровни логирования отличаются в зависимости от того, запущен проект локально, в режиме dev или в продакшене. По дефолту установлен локальный уровень.
Подробнее тут: pkg/logger/logger.go

## `Документация`
Swagger документация.

## Unit тесты
Тесты написаны на usecase с использованием моков (gomock):

Запуск тестов:
```sh
$ make test
```