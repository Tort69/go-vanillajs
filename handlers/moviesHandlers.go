package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Allusion/data"
	"Allusion/logger"
	"Allusion/models"
	"Allusion/token"
)

// Основная структура ответа API
type MovieResponse struct {
	Movie         models.Movie   `json:"movie"`
	RelatedMovies []models.Movie `json:"related_movies"`
}

type AllMoviesResponse struct {
	Movies   []data.MoviesPagination `json:"movies"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

type ActorResponse struct {
	Actor         models.Actor   `json:"actor"`
	RelatedMovies []models.Movie `json:"related_movies"`
}

type MovieHandler struct {
	storage data.MovieStorage
	logger  *logger.Logger
}

// Utility functions
func (h *MovieHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *MovieHandler) handleStorageError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		if err == data.ErrMovieNotFound {
			http.Error(w, context, http.StatusNotFound)
			return true
		}
		h.logger.Error(context, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}

func (h *MovieHandler) parseID(w http.ResponseWriter, idStr string) (int, bool) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid ID format", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func (h *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.storage.GetTopMovies()

	if h.handleStorageError(w, err, "Failed to get movies") {
		return
	}
	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served top movies")
	}
}

func (h *MovieHandler) GetRandomMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.storage.GetRandomMovies()
	if h.handleStorageError(w, err, "Failed to get movies") {
		return
	}
	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served random movies")
	}
}

func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {

	queryStr := r.URL.Query().Get("query")
	order := r.URL.Query().Get("order")
	genreStr := r.URL.Query().Get("genre")
	releaseYearStr := r.URL.Query().Get("releaseYear")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	var releaseYear *int
	page := 1
	pageSize := 50

	var query *string
	if queryStr == "" {
		query = nil
	} else {
		query = &queryStr
	}

	if releaseYearStr != "" {
		releaseYearStr, err := strconv.Atoi(releaseYearStr)
		if err != nil {
			http.Error(w, "Invalid page parameter: must be integer", http.StatusBadRequest)
			return
		}

		releaseYear = &releaseYearStr

	}

	if pageSizeStr != "" {
		parsedPageSizeStr, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			http.Error(w, "Invalid page parameter: must be integer", http.StatusBadRequest)
			return
		}

		pageSize = parsedPageSizeStr

	}

	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)

		if err != nil {
			http.Error(w, "Invalid page parameter: must be integer", http.StatusBadRequest)
			return
		}
		if parsedPage < 1 {
			http.Error(w, "Invalid page value: must be positive integer", http.StatusBadRequest)
			return
		}
		page = parsedPage
	}

	var genre *int
	if genreStr != "" {
		genreInt, ok := h.parseID(w, genreStr)
		if !ok {
			return
		}
		genre = &genreInt
	}

	var movies []data.MoviesPagination
	var err error
	movies, err = h.storage.SearchMoviesByName(query, order, genre, releaseYear, page, pageSize)

	if h.handleStorageError(w, err, "Failed to get movies") {
		return
	}

	response := AllMoviesResponse{
		Movies:   movies,
		Page:     page,
		PageSize: pageSize,
	}
	if h.writeJSONResponse(w, response) == nil {
		h.logger.Info("Successfully served movies")
	}
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {

	email, ok := token.ExtractJWTSecret(r, *h.logger)
	if !ok {
		h.logger.Info("User not authorized")
		email = ""
	}

	idStr := r.URL.Path[len("/api/movies/"):]
	id, ok := h.parseID(w, idStr)
	if !ok {
		return
	}

	movie, related_movies, err := h.storage.GetMovieByID(id, email)
	if h.handleStorageError(w, err, "Failed to get movie by ID") {
		return
	}

	response := MovieResponse{
		Movie:         movie,
		RelatedMovies: related_movies,
	}

	if h.writeJSONResponse(w, response) == nil {
		h.logger.Info("Successfully served movie with ID: " + idStr)
	}
}

func (h *MovieHandler) GetActor(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Path[len("/api/actor/"):]
	id, ok := h.parseID(w, idStr)
	if !ok {
		return
	}

	actor, movies, err := h.storage.GetMoviesActorById(id)
	if h.handleStorageError(w, err, "Failed to get actor by ID") {
		return
	}

	response := ActorResponse{
		Actor:         actor,
		RelatedMovies: movies,
	}

	if h.writeJSONResponse(w, response) == nil {
		h.logger.Info("Successfully served actor with ID: " + idStr)
	}
}

func (h *MovieHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.storage.GetAllGenres()
	if h.handleStorageError(w, err, "Failed to get genres") {
		return
	}
	if h.writeJSONResponse(w, genres) == nil {
		h.logger.Info("Successfully served genres")
	}
}

func NewMovieHandler(storage data.MovieStorage, log *logger.Logger) *MovieHandler {
	return &MovieHandler{
		storage: storage,
		logger:  log,
	}
}
