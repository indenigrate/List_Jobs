package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateJob(*Job) error
	DeleteJob(int) error
	GetJob() ([]*Job, error)
	GetJobByFilter(string, string, string) ([]*Job, error)
	GetJobByID(int) (*Job, error)
	UpdateJob(int, *Job) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(user, dbname, pass string) (*PostgresStore, error) {
	//save password in .env
	// connStr := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=disable", user, dbname, pass)
	connStr := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", user, pass, dbname)
	fmt.Printf("%v\n", connStr)
	// connStr := "user=postgres dbname=postgres password=12345 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	postgresStore := &PostgresStore{
		db: db,
	}
	if err = postgresStore.createJobTable(); err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}
	return postgresStore, nil
}
func (s *PostgresStore) Init() error {
	return nil
}

func (s *PostgresStore) createJobTable() error {
	//check if the table already exists
	var name sql.NullString
	query := "SELECT to_regclass('public.Job');"
	err := s.db.QueryRow(query).Scan(&name)
	if err != nil {
		return err
	}
	// fmt.Printf("name of table is %+v\n", name)

	//if not then create table
	query = `CREATE TABLE IF NOT EXISTS Job (
		id serial PRIMARY KEY NOT NULL,
		JobTitle   VARCHAR(100) NOT NULL,
		CompanyName VARCHAR(50) NOT NULL,
		Location    VARCHAR(50) NOT NULL,
		JobType     VARCHAR(50) NOT NULL,
		Description TEXT NOT NULL
	);`
	_, err = s.db.Exec(query)
	if err != nil {
		return err
	}

	if name.Valid {
		fmt.Println("Postgres Table exists:", name.String)
	} else {
		fmt.Println("Postgres Table does not exist, creating new table and initialising it with default test values")
		query = `INSERT INTO Job (JobTitle, CompanyName, Location, JobType, Description) VALUES
		('Software Engineer', 'Tech Solutions Inc.', 'New York, NY', 'Full-time', 'Develop scalable applications.'),
		('Product Manager', 'Innovate Corp', 'San Francisco, CA', 'Full-time', 'Lead product teams and strategies.'),
		('Data Analyst', 'Data Insights LLC', 'Remote', 'Part-time', 'Analyze data and create reports.'),
		('Marketing Specialist', 'Creative Agency', 'Remote', 'Contract', 'Implement marketing strategies.'),
		('UI/UX Designer', 'Design Hub', 'Chicago, IL', 'Full-time', 'Create user-friendly interfaces.'),
		('Systems Administrator', 'IT Services Co.', 'Remote', 'Full-time', 'Manage IT infrastructure and support.'),
		('Content Writer', 'Media Group', 'Remote', 'Freelance', 'Write engaging articles and content.'),
		('Sales Associate', 'Retail Corp', 'Seattle, WA', 'Part-time', 'Assist customers and drive sales.'),
		('DevOps Engineer', 'Cloud Solutions', 'Boston, MA', 'Full-time', 'Automate deployment processes.'),
		('Graphic Designer', 'Creative Studio', 'Portland, OR', 'Contract', 'Design visuals for marketing campaigns.');
		`
		_, err := s.db.Exec(query)
		if err != nil {
			return err
		}
	}
	//check if the table exists
	//if it does leave it as it is
	//if it does not initialise it with the following quer for testing purposes
	return nil
}

func (s *PostgresStore) CreateJob(job *Job) error {
	query := `INSERT INTO Job (JobTitle,CompanyName,Location,JobType,Description)
		VALUES (
        $1,
        $2,
        $3,
        $4,
        $5)
    RETURNING *;`
	_, err := s.db.Query(query, job.JobTitle, job.CompanyName, job.Location, job.JobType, job.Description)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return err
	}
	// fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostgresStore) DeleteJob(id int) error {
	_, err := s.db.Query("DELETE FROM Job WHERE id=$1", id)
	return err
}

func (s *PostgresStore) GetJobByID(id int) (*Job, error) {
	query := `SELECT * FROM Job WHERE id = $1;`
	resp, err := s.db.Query(query, id)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}
	defer resp.Close() // Ensure we close the rows after we're done

	for resp.Next() {
		// No rows found
		return scanJob(resp)
	}
	return nil, fmt.Errorf("account %d not found ", id)
}

func (s *PostgresStore) UpdateJob(id int, job *Job) error {
	if job.JobTitle != "" {
		_, err := s.db.Query("UPDATE Job SET JobTitle=$1 WHERE id=$2", job.JobTitle, id)
		if err != nil {
			return err
		}
	}
	if job.JobType != "" {
		_, err := s.db.Query("UPDATE Job SET JobType=$1 WHERE id=$2", job.JobType, id)
		if err != nil {
			return err
		}
	}
	if job.Location != "" {
		_, err := s.db.Query("UPDATE Job SET Location=$1 WHERE id=$2", job.Location, id)
		if err != nil {
			return err
		}
	}
	if job.Description != "" {
		_, err := s.db.Query("UPDATE Job SET Description=$1 WHERE id=$2", job.Description, id)
		if err != nil {
			return err
		}
	}
	if job.CompanyName != "" {
		_, err := s.db.Query("UPDATE Job SET CompanyName=$1 WHERE id=$2", job.CompanyName, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) GetJob() ([]*Job, error) {
	rows, err := s.db.Query("SELECT * FROM Job ORDER BY id ASC")
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}
	jobs := []*Job{}
	for rows.Next() {
		job := new(Job)
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s *PostgresStore) GetJobByFilter(JobTitle, location, JobType string) ([]*Job, error) {
	query := "SELECT * FROM Job WHERE 1=1"
	args := []interface{}{}
	// args := []string{}
	count := 1
	if JobTitle != "" {
		query += fmt.Sprintf(" AND JobTitle ILIKE $%d", count)
		args = append(args, "%"+JobTitle+"%")
		count += 1
	}
	if location != "" {
		// query += " AND location ILIKE $2"
		query += fmt.Sprintf(" AND Location ILIKE $%d", count)
		args = append(args, "%"+location+"%")
		count += 1
	}
	if JobType != "" {
		// query += " AND job_type = $3"
		query += fmt.Sprintf(" AND JobType ILIKE $%d", count)
		args = append(args, "%"+JobType+"%")
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}
	jobs := []*Job{}
	for rows.Next() {
		job := new(Job)
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func scanJob(rows *sql.Rows) (*Job, error) {
	job := new(Job)
	err := rows.Scan(
		&job.ID,
		&job.JobTitle,
		&job.CompanyName,
		&job.Location,
		&job.JobType,
		&job.Description)
	return job, err

}
