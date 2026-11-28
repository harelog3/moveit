// Package tmdb manages TMDB external API connection and relevant queries
package tmdb

type TmdbMovieSearchQuery struct {
	Page         int               `json:"page"`
	Results      []TmdbMovieSearch `json:"results"`
	TotalPages   int               `json:"total_pages"`
	TotalResults int               `json:"total_results"`
}

type TmdbMovieSearch struct {
	ID            int    `json:"id"`
	OriginalTitle string `json:"original_title"`
	ReleaseDate   string `json:"release_date"`
	PosterPath    string `json:"poster_path"`
}

type TmdbMovieDetail struct {
	Adult               bool                `json:"adult"`
	BackdropPath        string              `json:"backdrop_path"`
	BelongsToCollection *Collection         `json:"belongs_to_collection"` // null when the movie is not part of a collection
	Budget              int64               `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	ID                  int                 `json:"id"`
	IMDBID              string              `json:"imdb_id"`
	OriginCountry       []string            `json:"origin_country"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float64             `json:"popularity"`
	PosterPath          string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int64               `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Title               string              `json:"title"`
	Video               bool                `json:"video"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
}

type Collection struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type ProductionCountry struct {
	Iso3166_1 string `json:"iso_3166_1"`
	Name      string `json:"name"`
}

type SpokenLanguage struct {
	EnglishName string `json:"english_name"`
	Iso639_1    string `json:"iso_639_1"`
	Name        string `json:"name"`
}
