# List Jobs API

## Description

This is a Dockerized CRUD API built in Go, using PostgreSQL as the database.

## Prerequisites

- Docker
- Docker Compose

## Running the Application

1. Clone the repository:

   ```bash
   git clone https://github.com/indenigrate/list_jobs.git
   cd List_Jobs
2. Build and run:

   ```bash
   docker compose build
   docker compose up postgres -d
   docker compose up api -d

3. Access the API at http://localhost:8080/healthz.
## Basic Usage

### API Endpoints

Here are the available API endpoints for the List Jobs API:

- **POST `/login`**
  - **Description:** Authenticate a user and retrieve a JWT token.
  - **Request Body:** 
    ```json
    {
      "username": "admin",
      "password": "admin"
    }
    ```

- **GET `/healthz`**
  - **Description:** Check the health status of the API.
  - **Response:** Returns an empty json indicating the service is up.

- **GET `/jobs`**
  - **Description:** Retrieve a list of all jobs.
  - **Response:** Returns an array of job objects.

- **GET `/jobs/filter`**
  - **Description:** Retrieve a filtered list of jobs.
  - **Query Parameters:**
    - `job-title` (optional): The title of the job to filter.
    - `location` (optional): The location to filter jobs by.
    - `job-type` (optional): The type of job (e.g., Full-time, Part-time).
  - **Example:** `/jobs/filter?job_title=Engineer&location=NY&job_type=Full-time`

- **GET `/jobs/{id}`**
  - **Description:** Retrieve details of a specific job by its ID.
  - **URL Parameters:**
    - `id`: The ID of the job to retrieve.
  - **Response:** Returns the job object.

- **POST `/jobs`**
  - **Description:** Create a new job. Requires JWT authentication.
  - **Request Body:**
    ```json
    {
      "jobTitle": "Job Title",
      "companyName": "Company Name",
      "location": "Job Location",
      "jobType": "Job Type",
      "description": "Job Description"
    }
    ```

- **PUT `/jobs/{id}`**
  - **Description:** Update an existing job by its ID. Requires JWT authentication.
  - **URL Parameters:**
    - `id`: The ID of the job to update.
  - **Request Body:**
    ```json
    {
      "jobTitle": "Updated Job Title",
      "companyName": "Updated Company Name",
      "location": "Updated Job Location",
      "jobType": "Updated Job Type",
      "description": "Updated Job Description"
    }
    ```

- **DELETE `/jobs/{id}`**
  - **Description:** Delete a specific job by its ID. Requires JWT authentication.
  - **URL Parameters:**
    - `id`: The ID of the job to delete.




