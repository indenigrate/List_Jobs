package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	// "github.com/indenigrate/List_Jobs/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	// DB *database.Queries
	store Storage
}

func main() {
	//INITIALISATION
	//loading .env (environment variables)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Erro r loading .env file")
	}
	//retrieving PORT variable data
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	//initiating router
	router := chi.NewRouter()

	// initiating cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	//initialise storage
	store, err := NewPostgresStore("postgres", "postgres", "12345")
	if err != nil {
		log.Fatalf("unable to initialise storage: %v", err)
	}

	// fmt.Printf("%+v\n", store)
	apiCfg := apiConfig{store: store}
	//handle requests
	router.Get("/healthz", handlerReadiness)
	router.Get("/jobs", apiCfg.handlerListJob)
	router.Get("/jobs/{id}", apiCfg.handlerListJobByID)
	router.Post("/jobs", apiCfg.handlerCreateJob)
	router.Put("/jobs/{id}", apiCfg.handlerUpdateJob)
	router.Delete("/jobs/{id}", apiCfg.handlerDeleteJob)
	//initiate server properties
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %v\n", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server %v\n", err)
	}
}
