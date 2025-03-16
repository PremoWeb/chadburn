#!/bin/bash

# Simple Test Data Generator for Chadburn Metrics Dashboard
# This script adds jobs directly to the Chadburn configuration file

echo "Starting Simple Chadburn Test Data Generator"
echo "==========================================="

# Function to generate a random number between min and max
random_number() {
  min=$1
  max=$2
  echo $(( $min + RANDOM % ($max - $min + 1) ))
}

# Function to create a job-local entry
create_job_local() {
  job_id=$(random_number 1 1000)
  interval=$(random_number 5 30)
  
  echo "Creating job-local with ID $job_id (runs every $interval seconds)"
  
  # Append the job to the config file
  cat >> /etc/chadburn.conf << EOF

[job-local "test-job-$job_id"]
schedule = @every ${interval}s
command = echo "Test job $job_id running at \$(date)" && sleep \$(( RANDOM % 5 ))
EOF
}

# Function to create a job-local that occasionally fails
create_failing_job() {
  job_id=$(random_number 1001 2000)
  interval=$(random_number 10 40)
  
  echo "Creating occasionally failing job with ID $job_id (runs every $interval seconds)"
  
  # Append the job to the config file
  cat >> /etc/chadburn.conf << EOF

[job-local "failing-job-$job_id"]
schedule = @every ${interval}s
command = if [ \$(( RANDOM % 5 )) -eq 0 ]; then echo "Job $job_id failed" && exit 1; else echo "Job $job_id succeeded" && sleep \$(( RANDOM % 3 )); fi
EOF
}

# Function to create a long-running job
create_long_job() {
  job_id=$(random_number 2001 3000)
  interval=$(random_number 20 60)
  
  echo "Creating long-running job with ID $job_id (runs every $interval seconds)"
  
  # Append the job to the config file
  cat >> /etc/chadburn.conf << EOF

[job-local "long-job-$job_id"]
schedule = @every ${interval}s
command = echo "Starting long job $job_id" && sleep \$(( RANDOM % 10 + 5 )) && echo "Finished long job $job_id"
EOF
}

# Main loop
echo "Creating test jobs..."

# Create a mix of different job types
for i in {1..5}; do
  create_job_local
  sleep 2
done

for i in {1..3}; do
  create_failing_job
  sleep 2
done

for i in {1..2}; do
  create_long_job
  sleep 2
done

echo "Test jobs created. Waiting for metrics to be generated..."
echo "You should now see data in Prometheus and Grafana."
echo "Press Ctrl+C to stop this container (jobs will continue running)."

# Keep the container running
tail -f /dev/null 