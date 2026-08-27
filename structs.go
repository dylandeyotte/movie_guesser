package main

type ActorData struct {
	Page    int `json:"page"`
	Results []struct {
		Adult              bool         `json:"adult"`
		Gender             int          `json:"gender"`
		ID                 int          `json:"id"`
		KnownForDepartment string       `json:"known_for_department"`
		Name               string       `json:"name"`
		OriginalName       string       `json:"original_name"`
		Popularity         float64      `json:"popularity"`
		ProfilePath        string       `json:"profile_path"`
		KnownFor           []KnownFilms `json:"known_films"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type KnownFilms struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
	MediaType        string  `json:"media_type"`
	OriginalLanguage string  `json:"original_language"`
	GenreIds         []int   `json:"genre_ids"`
	Popularity       float64 `json:"popularity"`
	ReleaseDate      string  `json:"release_date"`
	Softcore         bool    `json:"softcore"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type FilmData struct {
	Adult               bool   `json:"adult"`
	BackdropPath        string `json:"backdrop_path"`
	BelongsToCollection struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		PosterPath   string `json:"poster_path"`
		BackdropPath string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`
	Budget int `json:"budget"`
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string   `json:"homepage"`
	ID                  int      `json:"id"`
	ImdbID              string   `json:"imdb_id"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalTitle       string   `json:"original_title"`
	Overview            string   `json:"overview"`
	Popularity          float64  `json:"popularity"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	Softcore        bool   `json:"softcore"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}
