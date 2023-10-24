package pagination

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

type (
	CustomKey string
)

const (
	NextPersonID CustomKey = "next_person_id"
)

// Pagination middleware is used to extract the next person id from the url query
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PersonID := r.URL.Query().Get(string(NextPersonID))
		intPersonID := 0

		var err error
		if PersonID != "" {
			intPersonID, err = strconv.Atoi(PersonID)
			if err != nil {
				_ = render.Render(w, r, ErrInvalidRequest(fmt.Errorf("couldn't read %s: %w", NextPersonID, err)))
				return
			}
		}

		ctx := context.WithValue(r.Context(), NextPersonID, intPersonID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status" example:"Resource not found."`                                         // user-level status message
	AppCode    int64  `json:"code,omitempty" example:"404"`                                                 // application-specific error code
	ErrorText  string `json:"error,omitempty" example:"The requested resource was not found on the server"` // application-level error message, for debugging
} // @name ErrorResponse

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns a structured http response for invalid requests
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
