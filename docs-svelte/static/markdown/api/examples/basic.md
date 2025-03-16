# Basic API Usage Examples

This page provides basic examples of how to use the Chadburn API with different programming languages and tools.

## Using cURL

### Authentication

```bash
# Using API key authentication
curl -X GET "http://your-chadburn-instance:8080/api/v1/jobs" \
  -H "X-API-Key: your-api-key"

# Using basic authentication
curl -X GET "http://your-chadburn-instance:8080/api/v1/jobs" \
  -u "username:password"

# Using JWT token authentication
# First, get a token
TOKEN=$(curl -X POST "http://your-chadburn-instance:8080/api/v1/auth/token" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}' | jq -r '.data.token')

# Then use the token
curl -X GET "http://your-chadburn-instance:8080/api/v1/jobs" \
  -H "Authorization: Bearer $TOKEN"
```

### Creating a Job

```bash
curl -X POST "http://your-chadburn-instance:8080/api/v1/jobs" \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-database",
    "schedule": "0 0 * * *",
    "command": "pg_dump -U postgres mydb > /backups/mydb.sql",
    "tags": ["backup", "database"]
  }'
```

### Listing Jobs

```bash
curl -X GET "http://your-chadburn-instance:8080/api/v1/jobs" \
  -H "X-API-Key: your-api-key"
```

### Getting a Specific Job

```bash
curl -X GET "http://your-chadburn-instance:8080/api/v1/jobs/job-123" \
  -H "X-API-Key: your-api-key"
```

### Updating a Job

```bash
curl -X PUT "http://your-chadburn-instance:8080/api/v1/jobs/job-123" \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-database-daily",
    "schedule": "0 0 * * *",
    "command": "pg_dump -U postgres mydb > /backups/mydb-$(date +%Y%m%d).sql",
    "tags": ["backup", "database", "daily"]
  }'
```

### Deleting a Job

```bash
curl -X DELETE "http://your-chadburn-instance:8080/api/v1/jobs/job-123" \
  -H "X-API-Key: your-api-key"
```

### Running a Job Immediately

```bash
curl -X POST "http://your-chadburn-instance:8080/api/v1/jobs/job-123/run" \
  -H "X-API-Key: your-api-key"
```

## Using Python

```python
import requests
import json

# API configuration
base_url = "http://your-chadburn-instance:8080/api/v1"
api_key = "your-api-key"
headers = {
    "X-API-Key": api_key,
    "Content-Type": "application/json"
}

# List all jobs
def list_jobs():
    response = requests.get(f"{base_url}/jobs", headers=headers)
    if response.status_code == 200:
        return response.json()["data"]
    else:
        print(f"Error: {response.status_code}")
        print(response.text)
        return None

# Get a specific job
def get_job(job_id):
    response = requests.get(f"{base_url}/jobs/{job_id}", headers=headers)
    if response.status_code == 200:
        return response.json()["data"]
    else:
        print(f"Error: {response.status_code}")
        print(response.text)
        return None

# Create a new job
def create_job(job_data):
    response = requests.post(f"{base_url}/jobs", headers=headers, json=job_data)
    if response.status_code == 201:
        return response.json()["data"]
    else:
        print(f"Error: {response.status_code}")
        print(response.text)
        return None

# Update a job
def update_job(job_id, job_data):
    response = requests.put(f"{base_url}/jobs/{job_id}", headers=headers, json=job_data)
    if response.status_code == 200:
        return response.json()["data"]
    else:
        print(f"Error: {response.status_code}")
        print(response.text)
        return None

# Delete a job
def delete_job(job_id):
    response = requests.delete(f"{base_url}/jobs/{job_id}", headers=headers)
    if response.status_code == 200:
        return response.json()["data"]
    else:
        print(f"Error: {response.status_code}")
        print(response.text)
        return None

# Run a job immediately
def run_job(job_id):
    response = requests.post(f"{base_url}/jobs/{job_id}/run", headers=headers)
    if response.status_code == 200:
        return response.json()["data"]
    else:
        print(f"Error: {response.status_code}")
        print(response.text)
        return None

# Example usage
if __name__ == "__main__":
    # Create a new job
    new_job = {
        "name": "backup-database",
        "schedule": "0 0 * * *",
        "command": "pg_dump -U postgres mydb > /backups/mydb.sql",
        "tags": ["backup", "database"]
    }
    
    job = create_job(new_job)
    if job:
        print(f"Created job: {job['id']}")
        
        # List all jobs
        jobs = list_jobs()
        if jobs:
            print(f"Found {len(jobs)} jobs")
            
        # Run the job immediately
        result = run_job(job['id'])
        if result:
            print(f"Job triggered: {result['execution_id']}")
```

## Using JavaScript

```javascript
// API configuration
const baseUrl = 'http://your-chadburn-instance:8080/api/v1';
const apiKey = 'your-api-key';
const headers = {
  'X-API-Key': apiKey,
  'Content-Type': 'application/json'
};

// List all jobs
async function listJobs() {
  try {
    const response = await fetch(`${baseUrl}/jobs`, {
      method: 'GET',
      headers
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error listing jobs:', error);
    return null;
  }
}

// Get a specific job
async function getJob(jobId) {
  try {
    const response = await fetch(`${baseUrl}/jobs/${jobId}`, {
      method: 'GET',
      headers
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error getting job:', error);
    return null;
  }
}

// Create a new job
async function createJob(jobData) {
  try {
    const response = await fetch(`${baseUrl}/jobs`, {
      method: 'POST',
      headers,
      body: JSON.stringify(jobData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error creating job:', error);
    return null;
  }
}

// Update a job
async function updateJob(jobId, jobData) {
  try {
    const response = await fetch(`${baseUrl}/jobs/${jobId}`, {
      method: 'PUT',
      headers,
      body: JSON.stringify(jobData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error updating job:', error);
    return null;
  }
}

// Delete a job
async function deleteJob(jobId) {
  try {
    const response = await fetch(`${baseUrl}/jobs/${jobId}`, {
      method: 'DELETE',
      headers
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error deleting job:', error);
    return null;
  }
}

// Run a job immediately
async function runJob(jobId) {
  try {
    const response = await fetch(`${baseUrl}/jobs/${jobId}/run`, {
      method: 'POST',
      headers
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error running job:', error);
    return null;
  }
}

// Example usage
async function main() {
  // Create a new job
  const newJob = {
    name: 'backup-database',
    schedule: '0 0 * * *',
    command: 'pg_dump -U postgres mydb > /backups/mydb.sql',
    tags: ['backup', 'database']
  };
  
  const job = await createJob(newJob);
  if (job) {
    console.log(`Created job: ${job.id}`);
    
    // List all jobs
    const jobs = await listJobs();
    if (jobs) {
      console.log(`Found ${jobs.length} jobs`);
    }
    
    // Run the job immediately
    const result = await runJob(job.id);
    if (result) {
      console.log(`Job triggered: ${result.execution_id}`);
    }
  }
}

main();
```

## Using Go

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// API configuration
const (
	baseURL = "http://your-chadburn-instance:8080/api/v1"
	apiKey  = "your-api-key"
)

// Job represents a Chadburn job
type Job struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name"`
	Schedule  string   `json:"schedule"`
	Command   string   `json:"command"`
	Tags      []string `json:"tags,omitempty"`
	Status    string   `json:"status,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
	NextRun   string   `json:"next_run,omitempty"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Data interface{} `json:"data"`
}

// Client is a simple Chadburn API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Chadburn API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// ListJobs retrieves all jobs
func (c *Client) ListJobs() ([]Job, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/jobs", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	// Convert the data to []Job
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	if err := json.Unmarshal(dataBytes, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetJob retrieves a specific job by ID
func (c *Client) GetJob(jobID string) (*Job, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/jobs/%s", c.baseURL, jobID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	// Convert the data to Job
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}

	var job Job
	if err := json.Unmarshal(dataBytes, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// CreateJob creates a new job
func (c *Client) CreateJob(job Job) (*Job, error) {
	jobBytes, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/jobs", c.baseURL), bytes.NewBuffer(jobBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	// Convert the data to Job
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}

	var createdJob Job
	if err := json.Unmarshal(dataBytes, &createdJob); err != nil {
		return nil, err
	}

	return &createdJob, nil
}

// Main function with example usage
func main() {
	client := NewClient(baseURL, apiKey)

	// Create a new job
	newJob := Job{
		Name:     "backup-database",
		Schedule: "0 0 * * *",
		Command:  "pg_dump -U postgres mydb > /backups/mydb.sql",
		Tags:     []string{"backup", "database"},
	}

	createdJob, err := client.CreateJob(newJob)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}

	fmt.Printf("Created job: %s\n", createdJob.ID)

	// List all jobs
	jobs, err := client.ListJobs()
	if err != nil {
		fmt.Printf("Error listing jobs: %v\n", err)
		return
	}

	fmt.Printf("Found %d jobs\n", len(jobs))

	// Get a specific job
	job, err := client.GetJob(createdJob.ID)
	if err != nil {
		fmt.Printf("Error getting job: %v\n", err)
		return
	}

	fmt.Printf("Job details: %+v\n", job)
} 