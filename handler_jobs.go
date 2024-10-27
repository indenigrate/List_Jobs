package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateJob(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		JobTitle    string `json:"job_title"`
		CompanyName string `json:"company_name"`
		Location    string `json:"location"`
		JobType     string `json:"job_type"`
		Description string `json:"description"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
	}
	job, err := apiCfg.DB.CreateUser(r.Context(), database.CreateJobParams{
		ID:          uuid.New(),
		JobTitle:    params.JobTitle,
		CompanyName: params.CompanyName,
		Location:    params.Location,
		JobType:     params.JobType,
		Description: params.Description,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create job:%v", err))
	}
	respondWithJSON(w, 201, job)
}

func (apiCfg *apiConfig) handlerListJob(w http.ResponseWriter, r *http.Request) {
}
func (apiCfg *apiConfig) handlerUpdateJob(w http.ResponseWriter, r *http.Request) {
}
func (apiCfg *apiConfig) handlerDeleteJob(w http.ResponseWriter, r *http.Request) {
}
