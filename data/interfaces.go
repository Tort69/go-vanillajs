package data

import "Allusion/models"

type MovieStorage interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
	GetMovieByID(id int) (models.Movie, error)
	SearchMoviesByName(name string, order string, genre *int) ([]models.Movie, error)
	GetAllGenres() ([]models.Genre, error)
}

type AccountStorage interface {
	Authenticate(string, string) (bool, error)
	Register(string, string, string) (bool, error)
	GetAccountDetails(string) (models.User, error)
	SaveCollection(models.User, int, string) (bool, error)
	DeleteCollection(models.User, int, string) (bool, error)
	DeleteAccount(string) (bool, error)
	VerifyEmail(string) (bool, string, error)
	ResendVerifyEmail(string) (bool, error)
	ResetPassword(string, string, string) (bool, error)
}
