package models

type (
	Movie struct {
		ID     int            `json:"id"`
		Title  string         `json:"title"`
		Genres []*MovieGenres `json:"genres"`
		Year   string         `json:"release_date"`
	}

	MovieGenres struct {
		Name string `json:"name"`
	}
)
