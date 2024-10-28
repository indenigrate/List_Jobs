package main

type CreateJobRequest struct {
	JobTitle    string `json:"job_title"`
	CompanyName string `json:"company_name"`
	Location    string `json:"location"`
	JobType     string `json:"job_type"`
	Description string `json:"description"`
}

type Job struct {
	ID          int    `json:"id"`
	JobTitle    string `json:"job_title"`
	CompanyName string `json:"company_name"`
	Location    string `json:"location"`
	JobType     string `json:"job_type"`
	Description string `json:"description"`
}

func NewJob(job_title, company_name, location, job_type, description string) *Job {
	return &Job{
		JobTitle:    job_title,
		CompanyName: company_name,
		Location:    location,
		JobType:     job_type,
		Description: description,
	}
}
