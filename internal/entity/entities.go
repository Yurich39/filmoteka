package entity

type Person struct {
	Id *int `db:"id" json:"id,omitempty"`
	Data
}

type Data struct {
	Name        *string `db:"name" json:"name,omitempty"`
	Surname     *string `db:"surname" json:"surname,omitempty"`
	Patronymic  *string `db:"patronymic" json:"patronymic,omitempty"`
	Age         *int    `db:"age" json:"age,omitempty"`
	Gender      *string `db:"gender" json:"gender,omitempty"`
	Nationality *string `db:"nationality" json:"nationality,omitempty"`
}
