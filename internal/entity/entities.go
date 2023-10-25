package entity

// DELETE по id
// GET по id или по Options
// PUT (это update) данных по Person - только по id с указанием новых данных в Data
// POST по данным в Data

type Person struct {
	Id
	Data
}

type Id struct {
	Id int `db:"id" json:"id"`
}

type Data struct {
	Name        *string `db:"name" json:"name,omitempty"`
	Surname     *string `db:"surname" json:"surname,omitempty"`
	Patronymic  *string `db:"patronymic" json:"patronymic,omitempty"`
	Age         *int    `db:"age" json:"age,omitempty"`
	Gender      *string `db:"gender" json:"gender,omitempty"`
	Nationality *string `db:"nationality" json:"nationality,omitempty"`
}
