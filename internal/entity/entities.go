package entity

type Person struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type UpdateRequest struct {
	Filters   map[string]interface{} `json:"filters,omitempty"`
	NewFields map[string]interface{} `json:"new_fields,omitempty"`
}

type DelRequest struct {
	Filters map[string]interface{} `json:"filters,omitempty"`
}

type Options struct {
	Where   map[string][]string
	OrderBy []string
	Order   []string
}
