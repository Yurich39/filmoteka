package webapi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

const op = "internal.usecase.webapi"

type EnrichWebAPI struct{}

func New() *EnrichWebAPI {
	return &EnrichWebAPI{}
}

type EnrichAgify struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func (t *EnrichWebAPI) EnrichAge(name string) (int, error) {

	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)

	r, err := http.Get(url)

	if err != nil {
		return -1, fmt.Errorf("%s: Agify API received issue: %w", op, err)
	}

	var req EnrichAgify
	err = render.DecodeJSON(r.Body, &req)
	if err != nil {
		return -1, fmt.Errorf("%s: Failed to decode Agify API request body: %w", op, err)
	}

	defer r.Body.Close()

	return req.Age, nil
}

type EnrichGenderize struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

func (t *EnrichWebAPI) EnrichGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	r, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("%s: Genderize API received issue: %w", op, err)
	}

	var req EnrichGenderize
	err = render.DecodeJSON(r.Body, &req)
	if err != nil {
		return "", fmt.Errorf("%s: Failed to decode Genderize API request body: %w", op, err)
	}

	defer r.Body.Close()

	return req.Gender, nil
}

type EnrichNationalize struct {
	Count   int                      `json:"count"`
	Name    string                   `json:"name"`
	Country []map[string]interface{} `json:"country"`
}

func (t *EnrichWebAPI) EnrichNationality(name string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	r, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("%s: Nationalize API received issue: %w", op, err)
	}

	var req EnrichNationalize
	err = render.DecodeJSON(r.Body, &req)
	if err != nil {
		return "", fmt.Errorf("%s: Failed to decode Genderize API request body: %w", op, err)
	}

	defer r.Body.Close()

	// Get string from interface{}
	res := ""
	nationality, ok := req.Country[0]["country_id"]
	if ok {
		res, _ = nationality.(string)
	}

	return res, nil
}
