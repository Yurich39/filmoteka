package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"filmoteka/config"
	"filmoteka/internal/controller/middleware/filter"
	"filmoteka/internal/controller/middleware/pagination"
	"filmoteka/internal/controller/middleware/sort"
	"filmoteka/internal/usecase"
	"filmoteka/pkg/logger"
)

func NewRouter(cfg *config.Config, router *chi.Mux, l logger.Interface, a usecase.Actor, m usecase.Movie, am usecase.ActorMovie) {
	// Middleware для общего использования
	commonMiddleware := chi.Chain(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.URLFormat,
	)

	// Middleware для авторизации администратора
	adminAuthMiddleware := middleware.BasicAuth("filmoteka", map[string]string{
		cfg.HTTPServer.User: cfg.HTTPServer.Pass,
	})

	actor := newActorHandler(a, l)
	movie := newMovieHandler(m, l)
	actor_movie := newActorMovieHandler(am, l)

	router.Route("/actor", func(r chi.Router) {
		r.Use(commonMiddleware.Handler)
		r.Get("/find/{id}", actor.find)
		r.With(adminAuthMiddleware).Post("/save", actor.save)
		r.With(adminAuthMiddleware).Put("/update", actor.update)
		r.With(adminAuthMiddleware).Delete("/delete/{id}", actor.delete)
	})

	router.Route("/actors", func(r chi.Router) {
		r.Route("/list", func(r chi.Router) {
			r.Use(commonMiddleware.Handler)
			r.With(filter.Middleware).Get("/", actor.list)
			r.With(filter.Middleware, pagination.Middleware).Get("/next", actor.next)
		})
	})

	router.Route("/movie", func(r chi.Router) {
		r.Use(commonMiddleware.Handler)
		r.Get("/find_by_id/{id}", movie.find)
		r.With(filter.Middleware).Get("/find/", movie.findMovie)
		r.With(adminAuthMiddleware).Post("/save", movie.save)
		r.With(adminAuthMiddleware).Put("/update", movie.update)
		r.With(adminAuthMiddleware).Delete("/delete/{id}", movie.delete)

	})

	router.Route("/movies", func(r chi.Router) {
		r.Route("/list", func(r chi.Router) {
			r.Use(commonMiddleware.Handler)
			r.With(filter.Middleware, sort.Middleware).Get("/", movie.list)
			r.With(filter.Middleware, sort.Middleware, pagination.Middleware).Get("/next", movie.next)
		})
	})

	router.Route("/actor_movie", func(r chi.Router) {
		r.Use(commonMiddleware.Handler)
		r.Get("/list", actor_movie.list)
		r.With(adminAuthMiddleware).Post("/save", actor_movie.save)
	})
}
