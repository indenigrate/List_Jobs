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
	store Storage
}

func main() {
	//INITIALISATION
	//loading .env (environment variables)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
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
	user := os.Getenv("PostgresUser")
	if user == "" {
		log.Fatal("PostgresUser is not found in the environment")
	}
	dbname := os.Getenv("PostgresDbname")
	if dbname == "" {
		log.Fatal("PostgresDbname is not found in the environment")
	}
	pass := os.Getenv("PostgresPass")
	if pass == "" {
		log.Fatal("PostgresPass is not found in the environment")
	}
	store, err := NewPostgresStore(user, dbname, pass)
	if err != nil {
		log.Fatalf("unable to initialise storage: %v", err)
	}

	// fmt.Printf("%+v\n", store)
	apiCfg := apiConfig{store: store}
	//handle requests
	router.Post("/login", apiCfg.handlerLogin)
	router.Get("/healthz", handlerReadiness)
	router.Get("/jobs", apiCfg.handlerListJob)
	router.Get("/jobs/filter", apiCfg.handlerListJobByFilter)
	//  /jobs/filter?job_title=JOB%20TITLE&location=JOB%20LOCATION&job_type=JOB%20TYPE
	router.Get("/jobs/{id}", apiCfg.handlerListJobByID)
	router.Post("/jobs", withJWTAuth(apiCfg.handlerCreateJob))
	router.Put("/jobs/{id}", withJWTAuth(apiCfg.handlerUpdateJob))
	router.Delete("/jobs/{id}", withJWTAuth(apiCfg.handlerDeleteJob))
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
