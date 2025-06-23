package data

import (
	"database/sql"
	"errors"
	"strconv"

	"Allusion/logger"
	"Allusion/models"

	_ "github.com/lib/pq"
)

type MovieRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewMovieRepository(db *sql.DB, log *logger.Logger) (*MovieRepository, error) {
	return &MovieRepository{
		db:     db,
		logger: log,
	}, nil
}

const defaultLimit = 20

func (r *MovieRepository) GetTopMovies() ([]models.Movie, error) {
	// Fetch movies
	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score,
		       popularity, language, poster_url, trailer_url
		FROM movies
		ORDER BY popularity DESC
		LIMIT 20
	`
	return r.getMovies(query)
}

func (r *MovieRepository) GetRandomMovies() ([]models.Movie, error) {
	// Fetch movies
	randomQuery := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score,
		       popularity, language, poster_url, trailer_url
		FROM movies
		ORDER BY random()
		LIMIT 20
	`
	return r.getMovies(randomQuery)
}

func (r *MovieRepository) GetAllMovies(offset int, pageSize int) ([]models.Movie, int, error) {
	// Fetch movies
	query := `
				SELECT
				id,
				tmdb_id,
				title,
				tagline,
				release_year,
				overview,
				score,
		    popularity,
				language,
				poster_url,
				trailer_url
        FROM movies
        ORDER BY release_year desc
        LIMIT 20 OFFSET 1;
`

	rows, err := r.db.Query(query)
	// , pageSize, offset
	if err != nil {
		r.logger.Error("Failed to query movies", err)
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan movie row", err)
			return nil, 0, err
		}
		movies = append(movies, m)
	}
	var totalCount int
	r.db.QueryRow("SELECT COUNT(*) FROM movies").Scan(&totalCount)

	return movies, totalCount, nil
}

func (r *MovieRepository) getMovies(query string) ([]models.Movie, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query movies", err)
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan movie row", err)
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func (r *MovieRepository) GetMovieByID(id int, email string) (models.Movie, []models.Movie, error) {
	var user models.User
	query := `
		SELECT id, name, email
		FROM users
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return models.Movie{}, []models.Movie{}, err
	}
	if err != nil {
		r.logger.Error("Failed to query user by email", err)
		return models.Movie{}, []models.Movie{}, err
	}
	// Fetch movie
	query = `
		SELECT m.id,
    m.tmdb_id,
    m.title,
    m.tagline,
    m.release_year,
    m.overview,
    m.score,
    m.popularity,
    m.language,
    m.poster_url,
    m.trailer_url,
		um.user_score,
    CONCAT(
        CASE
            WHEN MAX(CASE WHEN um.relation_type = 'favorite' THEN 1 ELSE 0 END) = 1 THEN 'In Favorite'
            ELSE 'Not Favorite'
        END,
        ', ',
        CASE
            WHEN MAX(CASE WHEN um.relation_type = 'watchlist' THEN 1 ELSE 0 END) = 1 THEN 'In Watchlist'
            ELSE 'Not Watchlist'
        END
    ) AS status
		FROM movies m
		LEFT JOIN user_movies um
				ON m.id = um.movie_id
				AND um.user_id = $1
		WHERE m.id = $2
		GROUP BY m.id, m.tmdb_id, m.title, m.tagline, m.release_year, m.overview, m.score, m.popularity, m.language, m.poster_url, m.trailer_url,um.user_score;
	`
	row := r.db.QueryRow(query, user.ID, id)

	var userScore sql.NullInt16
	var m models.Movie
	err = row.Scan(
		&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
		&m.Overview, &m.Score, &m.Popularity, &m.Language,
		&m.PosterURL, &m.TrailerURL, &userScore, &m.Status,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("Movie not found", ErrMovieNotFound)
		return models.Movie{}, []models.Movie{}, ErrMovieNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query movie by ID", err)
		return models.Movie{}, []models.Movie{}, err
	}
	if userScore.Valid {
		score := int(userScore.Int16)
		m.UserScore = &score
	} else {
		m.UserScore = nil
	}

	// Fetch related data
	if err := r.fetchMovieRelations(&m); err != nil {
		return models.Movie{}, []models.Movie{}, err
	}

	var related_movies []models.Movie
	related_movies, err = r.Related_MoviesById(&m)

	if err != nil {
		return models.Movie{}, []models.Movie{}, err
	}

	return m, related_movies, nil
}

func (r *MovieRepository) GetMoviesActorById(id int) (models.Actor, []models.Movie, error) {
	var actor models.Actor

	query := `
		SELECT id, first_name, last_name, image_url
		FROM actors
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&actor.ID,
		&actor.FirstName,
		&actor.LastName,
		&actor.ImageURL,
	)

	if err != nil {
		r.logger.Error("Actor not found", err)
		return models.Actor{}, []models.Movie{}, err
	}

	query = `SELECT  m.id,
						m.tmdb_id,
						m.title,
						m.tagline,
						m.release_year,
						m.overview,
						m.score,
						m.popularity,
						m.language,
						m.poster_url,
						m.trailer_url
						FROM movies m
						JOIN movie_cast mc ON m.id = mc.movie_id
						WHERE mc.actor_id = $1;`

	rows, err := r.db.Query(query, id)
	if err != nil {
		r.logger.Error("Failed to query movies", err)
		return models.Actor{}, []models.Movie{}, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan movie row", err)
			return models.Actor{}, []models.Movie{}, err
		}
		movies = append(movies, m)
	}
	return actor, movies, nil
}

func (r *MovieRepository) SearchMoviesByName(name *string, order string, genre *int, releaseYear *int, page int, pageSize int) ([]models.Movie, int, error) {
	orderBy := "release_year DESC"
	switch order {
	case "score":
		orderBy = "score DESC"
	case "name":
		orderBy = "title"
	case "popularity":
		orderBy = "popularity DESC"
	}

	// genreFilter := ""
	// if genre != nil {
	// 	genreFilter = strconv.Itoa(*genre)
	// }

	// releaseYearFilter := ""
	// if releaseYear != nil {
	// 	releaseYearFilter = strconv.Itoa(*releaseYear)
	// }
	// Fetch movies by name
	// query := `
	// 	SELECT id, tmdb_id, title, tagline, release_year, overview, score,
	// 	       popularity, language, poster_url, trailer_url, SELECT COUNT(*) as totalCount
	// 	FROM movies
	// 	WHERE (title ILIKE $1 OR overview ILIKE $1) ` + genreFilter + `
	// 	ORDER BY ` + orderBy + `
	// 	LIMIT $2 OFFSET $3;
	// `

	querySql := `WITH filtered_movies AS (
  SELECT
    id, tmdb_id, title, tagline, release_year,
    overview, score, popularity, language,
    poster_url, trailer_url,
    COUNT(*) OVER() AS total_count
  FROM movies
  WHERE
     ($1::INTEGER IS NULL OR EXISTS (
      SELECT 1 FROM movie_genres
      WHERE movie_id = movies.id AND genre_id = $1::INTEGER
    ))
    AND ($2::TEXT IS NULL OR title ILIKE '%%'|| $2::TEXT || '%%')
    AND ($3::INTEGER IS NULL OR release_year = $3::INTEGER)
)
SELECT *
FROM filtered_movies
ORDER BY $4::TEXT
LIMIT $5::INTEGER
OFFSET $6::INTEGER;`

	rows, err := r.db.Query(querySql, genre, name, releaseYear, orderBy, pageSize, page)
	if err != nil {
		r.logger.Error("Failed to search movies by name", err)
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	var totalCount int
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL, &totalCount,
		); err != nil {
			r.logger.Error("Failed to scan movie row", err)
			return nil, 0, err

		}
		movies = append(movies, m)
	}

	return movies, totalCount, nil
}

func (r *MovieRepository) GetAllGenres() ([]models.Genre, error) {
	query := `SELECT id, name FROM genres ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query all genres", err)
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			r.logger.Error("Failed to scan genre row", err)
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

// fetchMovieRelations fetches genres, actors, and keywords for a movie
func (r *MovieRepository) fetchMovieRelations(m *models.Movie) error {
	// Fetch genres
	genreQuery := `
		SELECT g.id, g.name
		FROM genres g
		JOIN movie_genres mg ON g.id = mg.genre_id
		WHERE mg.movie_id = $1
	`
	genreRows, err := r.db.Query(genreQuery, m.ID)
	if err != nil {
		r.logger.Error("Failed to query genres for movie "+strconv.Itoa(m.ID), err)
		return err
	}
	defer genreRows.Close()
	for genreRows.Next() {
		var g models.Genre
		if err := genreRows.Scan(&g.ID, &g.Name); err != nil {
			r.logger.Error("Failed to scan genre row", err)
			return err
		}
		m.Genres = append(m.Genres, g)
	}

	// Fetch actors
	actorQuery := `
		SELECT a.id, a.first_name, a.last_name, a.image_url
		FROM actors a
		JOIN movie_cast mc ON a.id = mc.actor_id
		WHERE mc.movie_id = $1
	`
	actorRows, err := r.db.Query(actorQuery, m.ID)
	if err != nil {
		r.logger.Error("Failed to query actors for movie "+strconv.Itoa(m.ID), err)
		return err
	}
	defer actorRows.Close()
	for actorRows.Next() {
		var a models.Actor
		if err := actorRows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.ImageURL); err != nil {
			r.logger.Error("Failed to scan actor row", err)
			return err
		}
		m.Casting = append(m.Casting, a)
	}

	// Fetch keywords
	keywordQuery := `
		SELECT k.word
		FROM keywords k
		JOIN movie_keywords mk ON k.id = mk.keyword_id
		WHERE mk.movie_id = $1
	`
	keywordRows, err := r.db.Query(keywordQuery, m.ID)
	if err != nil {
		r.logger.Error("Failed to query keywords for movie "+strconv.Itoa(m.ID), err)
		return err
	}
	defer keywordRows.Close()
	for keywordRows.Next() {
		var k string
		if err := keywordRows.Scan(&k); err != nil {
			r.logger.Error("Failed to scan keyword row", err)
			return err
		}
		m.Keywords = append(m.Keywords, k)
	}

	return nil
}

func (r *MovieRepository) Related_MoviesById(movies *models.Movie) ([]models.Movie, error) {

	query := `
		WITH target_genres AS (
				SELECT genre_id
				FROM movie_genres
				WHERE movie_id = $1
		)
		SELECT
				m.id,
				m.title,
				m.release_year,
				m.score,
				m.poster_url,
				COUNT(mg.genre_id) AS match_genre
		FROM movies m
		LEFT JOIN movie_genres mg
				ON m.id = mg.movie_id
				AND mg.genre_id IN (SELECT genre_id FROM target_genres)
		WHERE m.id != $1
		GROUP BY m.id, m.title, m.release_year, m.score, m.poster_url
		ORDER BY match_genre DESC, m.title
		LIMIT $2;
			`
	rows, err := r.db.Query(query, movies.ID, defaultLimit)
	if err != nil {
		r.logger.Error("Failed to query movies", err)
		return nil, err
	}
	defer rows.Close()

	var matchGenre int
	var related_movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.Title, &m.ReleaseYear,
			&m.Score, &m.PosterURL, &matchGenre,
		); err != nil {
			r.logger.Error("Failed to scan movie row", err)
			return nil, err
		}
		related_movies = append(related_movies, m)
	}

	return related_movies, nil
}

var (
	ErrMovieNotFound = errors.New("movie not found")
)
