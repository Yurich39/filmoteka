package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"people-finder/internal/controller/middleware/filter"
	"people-finder/internal/controller/middleware/pagination"
	"people-finder/internal/usecase"
	"people-finder/pkg/logger"
)

func NewRouter(router *chi.Mux, l logger.Interface, t usecase.Person) {
	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	person := newHandler(t, l)
	people := newHandler(t, l)

	router.Route("/person", func(r chi.Router) {

		// Handler #1
		r.Get("/find/{id}", person.find)

		// Handler #2
		r.Post("/save", person.save)

		// Handler #3
		r.Put("/update", person.update)

		// Handler #4
		r.Delete("/delete/{id}", person.delete)
	})

	router.Route("/people", func(r chi.Router) {

		r.Route("/list", func(r chi.Router) {

			// Middleware
			r.Use(filter.Middleware)

			// Handler #5
			r.Get("/", people.list)

			// Handler #6
			r.With(pagination.Middleware).Get("/next", people.next)
		})
	})
}
