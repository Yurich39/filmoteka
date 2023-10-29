package api

import "people-finder/internal/entity"

type Response struct {
	Status       string          `json:"status,omitempty"`
	Error        string          `json:"error,omitempty"`
	Person       *entity.Person  `json:"person,omitempty"`
	People       []entity.Person `json:"people,omitempty"`
	NextPersonID int             `json:"next_person_id,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
