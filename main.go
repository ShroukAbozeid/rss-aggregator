package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ShroukAbozeid/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/health", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users/me", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPIKey))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/me/feeds", apiCfg.middlewareAuth(apiCfg.handlerGetUserFeeds))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
