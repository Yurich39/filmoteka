package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"people-finder/internal/entity"
	"people-finder/internal/usecase"
	"people-finder/pkg/logger"
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

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if id == 0 {

		h.l.Info("No data available for id = 0")

		render.JSON(w, r,
			Response{
				Status: "Wrong id. Id should be > 0",
			})

		return
	}

	if err != nil {

		h.l.Debug("id parameter in URL is not integer or empty", h.l.Err(err))

		render.JSON(w, r, Error("Unable to retrieve id from URL"))

		return
	}

	res, err := h.t.Find(ctx, id)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to get data from DB"))

		return
	}

	render.JSON(w, r,
		Response{
			Status: StatusOk,
			Person: &res,
		})
}

func (h *handler) save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data entity.Data
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Data", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Data"))

		return
	}

	h.l.Info("request body decoded to entity.Data successfully", slog.Any("request", data))

	res, err := h.t.Save(ctx, data)

	if err != nil {
		h.l.Debug("Failed to save data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to save person data in DB"))

		return
	}

	render.JSON(w, r,
		Response{
			Status: StatusOk,
			Person: &res,
		})
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updates entity.Person
	err := render.DecodeJSON(r.Body, &updates)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Person", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Person"))

		return
	}

	h.l.Info("request body decoded to entity.Person successfully", slog.Any("request", updates))

	res, err := h.t.Update(ctx, updates)

	if err != nil {
		h.l.Debug("Failed to update data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to update person data in DB"))

		return
	}

	render.JSON(w, r,
		Response{
			Status: StatusOk,
			Person: &res,
		})
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {

		h.l.Debug("id parameter in URL is not integer or empty", h.l.Err(err))

		render.JSON(w, r, Error("Unable to retrieve id from URL"))

		return
	}

	h.l.Info("request body decoded to int successfully", slog.Any("request", id))

	if id == 0 {

		h.l.Info("No data available for id = 0")

		render.JSON(w, r,
			Response{
				Status: "Wrong id. Id should be > 0",
			})

		return
	}

	res, err := h.t.Delete(ctx, id)

	if err != nil {
		h.l.Debug("Failed to delete data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to delete data from DB"))

		return
	}

	render.JSON(w, r,
		Response{
			Status: StatusOk,
			Person: &res,
		})
}
