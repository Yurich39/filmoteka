package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"filmoteka/internal/entity"
	"filmoteka/internal/usecase"
	"filmoteka/pkg/logger"
)

type movieHandler struct {
	t usecase.Movie
	l logger.Interface
}

func newMovieHandler(t usecase.Movie, l logger.Interface) *movieHandler {
	return &movieHandler{t: t, l: l}
}

func (h *movieHandler) find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if id == 0 {

		h.l.Info("No data available for id = 0")

		render.JSON(w, r,
			MovieResponse{
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
		MovieResponse{
			Status: StatusOk,
			Movie:  &res,
		})
}

func (h *movieHandler) save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data entity.MovieData
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.MovieData", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.MovieData"))

		return
	}

	h.l.Info("request body decoded to entity.MovieData successfully", slog.Any("request", data))

	res, err := h.t.Save(ctx, data)

	if err != nil {
		h.l.Debug("Failed to save data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to save movie data in DB"))

		return
	}

	render.JSON(w, r,
		MovieResponse{
			Status: StatusOk,
			Movie:  &res,
		})
}

func (h *movieHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updates entity.Movie
	err := render.DecodeJSON(r.Body, &updates)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Movie", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Movie"))

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
		MovieResponse{
			Status: StatusOk,
			Movie:  &res,
		})
}

func (h *movieHandler) delete(w http.ResponseWriter, r *http.Request) {
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
			MovieResponse{
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
		MovieResponse{
			Status: StatusOk,
			Movie:  &res,
		})
}

func (h *movieHandler) findMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.t.FindMovie(ctx)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to get data from DB"))

		return
	}

	if len(res) == 0 {
		h.l.Debug("User input targeted does't target movie title or actor name")

		render.JSON(w, r, Error("No movies are in database with the speciefied input data"))

		return
	}

	render.JSON(w, r,
		MovieResponse{
			Status: StatusOk,
			Movies: res,
		})
}
