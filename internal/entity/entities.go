package entity

type Actor struct {
	Id *int `db:"id" json:"id,omitempty"`
	ActorData
}

type ActorData struct {
	Name        *string `db:"name" json:"name,omitempty"`
	Surname     *string `db:"surname" json:"surname,omitempty"`
	Patronymic  *string `db:"patronymic" json:"patronymic,omitempty"`
	Gender      *string `db:"gender" json:"gender,omitempty"`
	DateOfBirth *string `db:"date_of_birth" json:"date_of_birth,omitempty"`
}

type Movie struct {
	Id *int `db:"id" json:"id,omitempty"`
	MovieData
}

type MovieData struct {
	Title       *string `db:"title" json:"title,omitempty"`
	Description *string `db:"description" json:"description,omitempty"`
	ReleaseDate *string `db:"release_date" json:"release_date,omitempty"`
	Rating      *int    `db:"rating" json:"rating,omitempty"`
}

type ActorMovie struct {
	Actor_id *int `db:"actor_id" json:"actor_id,omitempty"`
	Movie_id *int `db:"movie_id" json:"movie_id,omitempty"`
}

type ActorMovieData struct {
	ActorID      int    `db:"actor_id"`
	ActorName    string `db:"actor_name"`
	ActorSurname string `db:"actor_surname"`
	MovieID      int    `db:"movie_id"`
	MovieTitle   string `db:"movie_title"`
}

type MoviesOfActor struct {
	ActorID      int      `db:"actor_id" json:"actor_id,omitempty"`
	ActorName    string   `db:"actor_name" json:"actor_name,omitempty"`
	ActorSurname string   `db:"actor_surname" json:"actor_surname,omitempty"`
	Movies       []string `db:"movies" json:"movies,omitempty"`
}
