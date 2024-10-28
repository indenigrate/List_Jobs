package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (apiCfg *apiConfig) handlerCreateJob(w http.ResponseWriter, r *http.Request) {
	job, err := readJobInput(r)
	if err != nil {
		respondWithError(w, 400, err.Error())
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
