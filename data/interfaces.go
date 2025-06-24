package data

import "Allusion/models"

type MovieStorage interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
	GetAllMovies(int, int) ([]models.Movie, int, error)
	GetMovieByID(id int, email string) (models.Movie, []models.Movie, error)
	SearchMoviesByName(name *string, order string, genre *int, releaseYear *int, page int, pageSize int) ([]MoviesPagination,  error)
	GetAllGenres() ([]models.Genre, error)
	GetMoviesActorById(id int) (models.Actor, []models.Movie, error)
}

type AccountStorage interface {
	Authenticate(string, string) (bool, error)
	Register(string, string, string) (bool, error)
	GetAccountDetails(string) (models.User, error)
	SaveCollection(models.User, int, string, *int) (bool, error)
	DeleteCollection(models.User, int, string) (bool, error)
	DeleteAccount(string) (bool, error)
	VerifyEmail(string) (bool, string, error)
	ResendVerifyEmail(string) (bool, error)
	ResetPassword(string, string, string) (bool, error)
}
