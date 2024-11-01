package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	jwt "github.com/golang-jwt/jwt/v5"
)

func (apiCfg *apiConfig) handlerCreateJob(w http.ResponseWriter, r *http.Request) {
	// Get the token from the request context
	token := r.Context().Value("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	// Check if the user role is admin
	if claims["role"] != "admin" {
		respondWithError(w, http.StatusForbidden, "access denied: admin only")
		return
	}

	job, err := readJobInput(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	if err := apiCfg.store.CreateJob(job); err != nil {
		respondWithError(w, 400, "something went wrong in handlerCreateJob")
		return
	}
	respondWithJSON(w, 201, job)
}

func (apiCfg *apiConfig) handlerListJob(w http.ResponseWriter, r *http.Request) {
	jobs, err := apiCfg.store.GetJob()
	if err != nil {
		respondWithError(w, 400, "Error in handlerListJob")
	}
	respondWithJSON(w, 200, jobs)
}

func (apiCfg *apiConfig) handlerListJobByID(w http.ResponseWriter, r *http.Request) {
	//get id
	id, err := getID(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	//call getjobbyid func in storage
	job, err := apiCfg.store.GetJobByID(id)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error in handlerListJobByID: %v", err))
		return
	}
	// fmt.Printf("%v %+v\n", intId, job)
	respondWithJSON(w, 200, job)
}

func (apiCfg *apiConfig) handlerUpdateJob(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	// Check if the user role is admin
	if claims["role"] != "admin" {
		respondWithError(w, http.StatusForbidden, "access denied: admin only")
		return
	}
	//read id
	id, err := getID(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	//read updated job data
	job, err := readJobInput(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	err = apiCfg.store.UpdateJob(id, job)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}
	respondWithJSON(w, 201, fmt.Sprintf("Job %d updated succesfully", id))
}

func (apiCfg *apiConfig) handlerDeleteJob(w http.ResponseWriter, r *http.Request) {

	token := r.Context().Value("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	// Check if the user role is admin
	if claims["role"] != "admin" {
		respondWithError(w, http.StatusForbidden, "access denied: admin only")
		return
	}
	//get id
	id, err := getID(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	//call delete function in storage
	if err := apiCfg.store.DeleteJob(id); err != nil {
		respondWithError(w, 400, fmt.Sprintf("error: %v", err))
		return
	}
	respondWithJSON(w, 200, fmt.Sprintf("Index %d deleted succesfully", id))

}

func (apiCfg *apiConfig) handlerListJobByFilter(w http.ResponseWriter, r *http.Request) {
	JobTitle := r.URL.Query().Get("job_title")
	location := r.URL.Query().Get("location")
	jobType := r.URL.Query().Get("job_type")

	jobs, err := apiCfg.store.GetJobByFilter(JobTitle, location, jobType)
	if err != nil {
		respondWithError(w, 400, "Error in handlerListJobByFilter")
	}
	respondWithJSON(w, 200, jobs)
}

func getID(r *http.Request) (int, error) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return intId, fmt.Errorf("invalid id given: %v", err)
	}
	return intId, nil
}

func readJobInput(r *http.Request) (*Job, error) {
	CreateJobRequest := new(CreateJobRequest)
	if err := json.NewDecoder(r.Body).Decode(CreateJobRequest); err != nil {
		return nil, err
	}
	job := NewJob(CreateJobRequest.JobTitle, CreateJobRequest.CompanyName, CreateJobRequest.Location, CreateJobRequest.JobType, CreateJobRequest.Description)
	return job, nil
}

// create and return jwt(json web Token) when admin logs in
// default username=admin
// default password=admin
func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Validate user credentials (this is a placeholder)
	userRole := "user" // Default role
	if loginRequest.Username == "admin" && loginRequest.Password == "admin" {
		userRole = "admin"
	} else if loginRequest.Username != "user" {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Create JWT
	tokenString, err := createJWT(userRole)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}
