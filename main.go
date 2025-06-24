package main

import (
	"Allusion/data"
	"Allusion/handlers"
	"Allusion/logger"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
)

func main() {

	logInstance := initializeLogger()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	// Проверка подключения
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Connection failed:%v", err)

	}
	log.Printf("Connected:%v", pong)

	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or failed to load: %v", err)
	}

	// dbConnStr := os.Getenv("DATABASE_URL")
	// if dbConnStr == "" {
	// 	log.Fatalf("DATABASE_URL not set in environment")
	// }
	// db, err := sql.Open("postgres", dbConnStr)
	connConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	connConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	db, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	accountRepo, err := data.NewAccountRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize account repository: %v", err)
	}
	accountHandler := handlers.NewAccountHandler(accountRepo, logInstance, rdb)

	movieRepo, err := data.NewMovieRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize movie repository: %v", err)
	}
	movieHandler := handlers.NewMovieHandler(movieRepo, logInstance)

	http.HandleFunc("/api/movies/random/", movieHandler.GetRandomMovies)
	http.HandleFunc("/api/movies/top/", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/search/", movieHandler.SearchMovies)
	http.HandleFunc("/api/movies/", movieHandler.GetMovie)
	// http.HandleFunc("/api/movies/all-movies", movieHandler.GetAllMovies)
	// http.HandleFunc("/api/movies/all-movies/", movieHandler.HandlerGetAllMovies)
	http.HandleFunc("/api/actor/", movieHandler.GetActor)
	http.HandleFunc("/api/genres/", movieHandler.GetGenres)
	http.HandleFunc("/api/account/register/", accountHandler.Register)
	http.HandleFunc("/api/account/authenticate/", accountHandler.Authenticate)
	http.HandleFunc("/api/account/verify/",
		accountHandler.VerifyByEmail)
	http.Handle("/api/account/resendVerifyEmail/",
		accountHandler.RateLimitMiddleware(http.HandlerFunc(accountHandler.HandlerResendVerifyEmail)))

	http.Handle("/api/account/resetPassword/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.HandlerResetPassword)))

	http.Handle("/api/account/favorites/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetFavorites)))

	http.Handle("/api/account/deleteAccount/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.DeleteAccountHandler)))

	http.Handle("/api/account/watchlist/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetWatchlist)))

	http.Handle("/api/account/save-to-collection/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.SaveToCollection)))

	http.Handle("/api/account/delete-to-collection/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.DeleteToCollection)))

	// Handler catch-all
	catchAllHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	}
	http.HandleFunc("/movies", catchAllHandler)
	http.HandleFunc("/movies/", catchAllHandler)
	http.HandleFunc("/actor/", catchAllHandler)
	http.HandleFunc("/account/", catchAllHandler)
	http.HandleFunc("/account/verify", catchAllHandler)

	http.Handle("/", http.FileServer(http.Dir("public")))

	const addr = ":8080"
	logInstance.Info("Server starting on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logInstance.Error("Server failed to start", err)
		log.Fatalf("Server failed: %v", err)
	}
}

func initializeLogger() *logger.Logger {
	logInstance, err := logger.NewLogger("movie-service.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logInstance.Close()
	return logInstance
}
