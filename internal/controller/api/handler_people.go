package api

import (
	"net/http"

	"github.com/go-chi/render"
)

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.t.List(ctx)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to get data from DB"))

		return
	}

	id := res[len(res)-1].Id.Id

	render.JSON(w, r,
		Response{
			Status:       StatusOk,
			People:       res,
			NextPersonID: id + 1,
		})
}

// Получаем пагинацию
func (h *handler) next(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	res, err := h.t.Next(ctx)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to get data from DB"))

		return
	}

	id := res[len(res)-1].Id.Id

	render.JSON(w, r,
		Response{
			Status:       StatusOk,
			People:       res,
			NextPersonID: id + 1,
		})

}
