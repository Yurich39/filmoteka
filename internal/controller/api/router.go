package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"people-finder/internal/usecase"
	"people-finder/pkg/logger"

	"people-finder/internal/controller/middleware/filter"
	"people-finder/internal/controller/middleware/pagination"
	"people-finder/internal/controller/middleware/sort"
)

func NewRouter(router *chi.Mux, l logger.Interface, t usecase.Person) {
	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	handler := newHandler(t, l)

	router.Route("/person", func(r chi.Router) {

		// Handler #1
		r.Post("/save", handler.save)

		// Handler #2
		r.Put("/update", handler.update)

		// Handler #3
		r.Delete("/delete", handler.delete)
	})

	router.Route("/people", func(r chi.Router) {
		// Middleware
		r.Use(filter.Middleware)
		r.Use(sort.Middleware)

		// Handler #4
		r.With(pagination.Middleware).Get("/next", handler.listPeople)

		r.Route("/find", func(r chi.Router) {
			// Handler #5
			r.Get("/", handler.find)
		})
	})
}
