package api

import (
	"net/http"
	"strconv"

	"people-finder/internal/usecase"
	"people-finder/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"people-finder/internal/entity"
)

type handler struct {
	t usecase.Person
	l logger.Interface
}

func newHandler(t usecase.Person, l logger.Interface) *handler {
	return &handler{t: t, l: l}
}

func (h *handler) find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var id entity.Id
	val, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		h.l.Debug("id parameter in URL is not integer or empty", h.l.Err(err))

		render.JSON(w, r, Error("Unable to retrieve id from URL"))

		return
	}

	id.Id = val

	res, err := h.t.Find(ctx, id)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to get data from DB"))

		return
	}

	render.JSON(w, r,
		Response{
			Status: StatusOk,
			People: res,
		})
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
	data := Enrich(h, *person.Data.Name)

	h.l.Info("API enriched data received successfully:", slog.Any("request", data))

	person.Data.Age = &data.Age
	person.Data.Gender = &data.Gender
	person.Data.Nationality = &data.Nationality

	// Save data in db
	res, err := h.t.Save(ctx, person)

	if err != nil {
		h.l.Debug("Failed to save data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to save person data in DB"))

		return
	}

	render.JSON(w, r,
		Response{
			Status: StatusOk,
			People: res,
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

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode JSON
	var updates entity.Person
	err := render.DecodeJSON(r.Body, &updates)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Person", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Person"))

		return
	}

	h.l.Info("request body decoded to entity.Person successfully", slog.Any("request", updates))

	// Validate data

	// Update data in db
	res, err := h.t.Update(ctx, updates)

	if err != nil {
		h.l.Debug("Failed to update data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to update person data in DB"))

		return
	}

	render.JSON(w, r,
		Response{
			Status: StatusOk,
			People: res,
		})
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode JSON
	var id entity.Id
	err := render.DecodeJSON(r.Body, &id)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Id", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Id"))

		return
	}

	h.l.Info("request body decoded to entity.Id successfully", slog.Any("request", id))

	// Validate data

	// Delete data from db
	res, err := h.t.Delete(ctx, id)

	if err != nil {
		h.l.Debug("Failed to delete data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to delete data from DB"))

		return
	}

	render.JSON(w, r,
		Response{
			Status: StatusOk,
			People: res,
		})
}
