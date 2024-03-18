package api

import "filmoteka/internal/entity"

type ActorResponse struct {
	Status      string         `json:"status,omitempty"`
	Actor       *entity.Actor  `json:"actor,omitempty"`
	Actors      []entity.Actor `json:"actors,omitempty"`
	NextActorID int            `json:"next_actor_id,omitempty"`
}

type MovieResponse struct {
	Status      string         `json:"status,omitempty"`
	Movie       *entity.Movie  `json:"movie,omitempty"`
	Movies      []entity.Movie `json:"movies,omitempty"`
	NextMovieID int            `json:"next_movie_id,omitempty"`
}

type ActorMovieResponse struct {
	Status     string             `json:"status,omitempty"`
	ActorMovie *entity.ActorMovie `json:"actor_movie,omitempty"`
}

type ActorsMoviesResponse struct {
	Status string                 `json:"status,omitempty"`
	Data   []entity.MoviesOfActor `json:"movies_of_actor,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

type customError struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Error(msg string) customError {
	return customError{
		Status: StatusError,
		Error:  msg,
	}
}
