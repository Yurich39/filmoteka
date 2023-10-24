package api

import (
	"net/http"

	"people-finder/internal/usecase"
	"people-finder/pkg/logger"

	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"people-finder/internal/controller/middleware/filter"
	"people-finder/internal/controller/middleware/pagination"
	"people-finder/internal/controller/middleware/sort"
	"people-finder/internal/entity"
)

type handler struct {
	t usecase.Person
	l logger.Interface
}

func newHandler(t usecase.Person, l logger.Interface) *handler {
	return &handler{t: t, l: l}
}

type EnrichData struct {
	Age         int
	Gender      string
	Nationality string
}

func (h *handler) save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode JSON
	var person entity.Person
	err := render.DecodeJSON(r.Body, &person)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Person", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Person"))

		return
	}

	h.l.Info("request body decoded to entity.Person successfully", slog.Any("request", person))

	// Validate data

	// Enrich data
	data := Enrich(h, person.Name)

	h.l.Info("API enriched data received successfully:", slog.Any("request", data))

	person.Age = data.Age
	person.Gender = data.Gender
	person.Nationality = data.Nationality

	// Save data in db
	res, err := h.t.Save(ctx, person)

	if err != nil {
		h.l.Debug("Failed to save data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to save person data in DB"))

		return
	}

	render.JSON(w, r, Response{
		Status: StatusOk,
		People: res,
	})

}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode JSON
	var updates entity.UpdateRequest
	err := render.DecodeJSON(r.Body, &updates)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.UpdateRequest", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.UpdateRequest"))

		return
	}

	h.l.Info("request body decoded to entity.UpdateRequest successfully", slog.Any("request", updates))

	// Validate data

	// Update data in db
	res, err := h.t.Update(ctx, updates)

	if err != nil {
		h.l.Debug("Failed to update data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to update person data in DB"))

		return
	}

	render.JSON(w, r, Response{
		Status: StatusOk,
		People: res,
	})
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode JSON
	var deleter entity.DelRequest
	err := render.DecodeJSON(r.Body, &deleter)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.DelRequest", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.DelRequest"))

		return
	}

	h.l.Info("request body decoded to entity.DelRequest successfully", slog.Any("request", deleter))

	// Validate data

	// Delete data from db
	res, err := h.t.Delete(ctx, deleter)

	if err != nil {
		h.l.Debug("Failed to delete data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to delete data from DB"))

		return
	}

	render.JSON(w, r, Response{
		Status: StatusOk,
		People: res,
	})
}

func (h *handler) find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Достаем url параметры (ключ/значение) из http запроса

	// Делаем type assertion
	var sortOptions sort.Options

	switch val := ctx.Value(sort.SortOptionsContextKey).(type) {
	case nil:
		sortOptions = sort.Options{}
	case sort.Options:
		sortOptions = sort.Options(val)
	default:
		sortOptions := sort.Options{}
		h.l.Debug("expected string-keyed map or string, not %T", sortOptions)
	}

	var filterOptions filter.Options
	switch val := r.Context().Value(filter.FilterOptionsContextKey).(type) {
	case nil:
		filterOptions = filter.Options{}
	case filter.Options:
		filterOptions = filter.Options(val)
	default:
		filterOptions = filter.Options{}
		h.l.Debug("expected string-keyed map or string, not %T", filterOptions)
	}

	// Преобразовываем параметры в структуру, относящуюся к БД
	getter := entity.Options{
		Where:   filterOptions.Where,
		OrderBy: sortOptions.OrderBy,
		Order:   sortOptions.Order,
	}

	res, err := h.t.Find(ctx, getter)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to get data from DB"))

		return
	}

	id := res[len(res)-1].ID

	render.JSON(w, r, Response{
		Status:       StatusOk,
		People:       res,
		NextPersonID: id + 1,
	})
}

// Получаем пагинацию
func (h *handler) listPeople(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Достаем url параметры (ключ/значение) из http запроса

	// Делаем type assertion
	var sortOptions sort.Options

	switch val := ctx.Value(sort.SortOptionsContextKey).(type) {
	case nil:
		sortOptions = sort.Options{}
	case sort.Options:
		sortOptions = sort.Options(val)
	default:
		sortOptions := sort.Options{}
		h.l.Debug("expected string-keyed map or string, not %T", sortOptions)
	}

	var filterOptions filter.Options
	switch val := r.Context().Value(filter.FilterOptionsContextKey).(type) {
	case nil:
		filterOptions = filter.Options{}
	case filter.Options:
		filterOptions = filter.Options(val)
	default:
		filterOptions = filter.Options{}
		h.l.Debug("expected string-keyed map or string, not %T", filterOptions)
	}

	// Преобразовываем параметры в структуру, относящуюся к БД
	getter := entity.Options{
		Where:   filterOptions.Where,
		OrderBy: sortOptions.OrderBy,
		Order:   sortOptions.Order,
	}

	personID := r.Context().Value(pagination.NextPersonID).(int)

	res, err := h.t.ListPeople(ctx, getter, personID)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to get data from DB"))

		return
	}

	id := res[len(res)-1].ID

	render.JSON(w, r, Response{
		Status:       StatusOk,
		People:       res,
		NextPersonID: id + 1,
	})

}

func Enrich(h *handler, name string) EnrichData {
	data := EnrichData{}

	// Request Age
	age, err := h.t.GetAge(name)

	if err != nil {
		h.l.Debug("API fail:", err)
	} else {
		data.Age = age
	}

	// Request Gender
	gender, err := h.t.GetGender(name)

	if err != nil {
		h.l.Debug("API fail:", err)
	} else {
		data.Gender = gender
	}

	// Request Nationality
	nationality, err := h.t.GetNationality(name)

	if err != nil {
		h.l.Debug("API fail:", err)
	} else {
		data.Nationality = nationality
	}

	return data
}
