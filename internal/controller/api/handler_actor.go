package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lib/pq"
	"golang.org/x/exp/slog"

	"filmoteka/internal/entity"
	"filmoteka/internal/usecase"
	"filmoteka/pkg/logger"
)

type actorHandler struct {
	t usecase.Actor
	l logger.Interface
}

func newActorHandler(t usecase.Actor, l logger.Interface) *actorHandler {
	return &actorHandler{t: t, l: l}
}

func (h *actorHandler) find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if id == 0 {

		h.l.Info("No data available for id = 0")

		render.JSON(w, r,
			ActorResponse{
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

		render.JSON(w, r, Error(fmt.Sprintf("Database has NO actor with id = %d", id)))

		return
	}

	render.JSON(w, r,
		ActorResponse{
			Status: StatusOk,
			Actor:  &res,
		})
}

func (h *actorHandler) save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data entity.ActorData
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.ActorData", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.ActorData"))

		return
	}

	h.l.Info("request body decoded to entity.ActorData successfully", slog.Any("request", data))

	res, err := h.t.Save(ctx, data)

	if err != nil {
		h.l.Debug("Failed to save data in DB", h.l.Err(err))

		if err == sql.ErrNoRows {
			render.JSON(w, r, Error("actor already exists"))
			return
		} else {
			switch err := err.(type) {
			case *pq.Error:
				render.JSON(w, r, Error("provided data is invalid"))
				return
			default:
				render.JSON(w, r, err.Error())
				return
			}
		}
	}

	render.JSON(w, r,
		ActorResponse{
			Status: StatusOk,
			Actor:  &res,
		})
}

func (h *actorHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updates entity.Actor
	err := render.DecodeJSON(r.Body, &updates)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Actor", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Actor"))

		return
	}

	h.l.Info("request body decoded to entity.Actor successfully", slog.Any("request", updates))

	res, err := h.t.Update(ctx, updates)

	if err != nil {
		h.l.Debug("Failed to update data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to update actor data in DB"))

		return
	}

	render.JSON(w, r,
		ActorResponse{
			Status: StatusOk,
			Actor:  &res,
		})
}

func (h *actorHandler) delete(w http.ResponseWriter, r *http.Request) {
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
			ActorResponse{
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
		ActorResponse{
			Status: StatusOk,
			Actor:  &res,
		})
}
