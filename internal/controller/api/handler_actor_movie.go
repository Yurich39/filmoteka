package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/lib/pq"
	"golang.org/x/exp/slog"

	"filmoteka/internal/entity"
	"filmoteka/internal/usecase"
	"filmoteka/pkg/logger"
)

type actorMovieHandler struct {
	t usecase.ActorMovie
	l logger.Interface
}

func newActorMovieHandler(t usecase.ActorMovie, l logger.Interface) *actorMovieHandler {
	return &actorMovieHandler{t: t, l: l}
}

func (h *actorMovieHandler) save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data entity.ActorMovie
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.ActorMovie", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.ActorMovie"))

		return
	}

	h.l.Info("request body decoded to entity.ActorMovie successfully", slog.Any("request", data))

	err = h.t.Save(ctx, data)

	if err != nil {
		h.l.Debug("Failed to save data in DB", h.l.Err(err))

		switch err := err.(type) {
		case *pq.Error:
			render.JSON(w, r, Error("provided data is invalid or actor_id to movie_id assignment already exists"))
			return
		default:
			render.JSON(w, r, err.Error())
			return
		}
	}

	render.JSON(w, r,
		ActorMovieResponse{
			Status:     StatusOk,
			ActorMovie: &data,
		})
}

func (h *actorMovieHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.t.List(ctx)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to get data from DB"))

		return
	}

	if len(res) == 0 {
		h.l.Info("No data")

		render.JSON(w, r, customError{
			Status: StatusOk,
			Error:  "No data",
		})

		return
	}

	//Представим данные из res в требуемый формат
	newData := ConvertToMoviesOfActor(res)

	render.JSON(w, r,
		ActorsMoviesResponse{
			Status: StatusOk,
			Data:   newData,
		})
}

type ActorCridentials struct {
	name    string
	surname string
}

func ConvertToMoviesOfActor(data []entity.ActorMovieData) []entity.MoviesOfActor {
	moviesByActor := make(map[int][]string)
	ActorData := make(map[int]ActorCridentials)

	for _, d := range data {
		moviesByActor[d.ActorID] = append(moviesByActor[d.ActorID], d.MovieTitle)
		ActorData[d.ActorID] = ActorCridentials{name: d.ActorName, surname: d.ActorSurname}
	}

	var moviesOfActors []entity.MoviesOfActor

	for actorID, movies := range moviesByActor {

		movieOfActor := entity.MoviesOfActor{
			ActorID:      actorID,
			ActorName:    ActorData[actorID].name,
			ActorSurname: ActorData[actorID].surname,
			Movies:       movies,
		}
		moviesOfActors = append(moviesOfActors, movieOfActor)
	}

	return moviesOfActors
}
