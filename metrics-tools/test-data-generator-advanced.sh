#!/bin/bash

# Advanced Test Data Generator for Chadburn Metrics Dashboard
# This script creates various types of Chadburn jobs and generates diverse metrics

echo "Starting Advanced Chadburn Test Data Generator"
echo "=============================================="

# Function to generate a random number between min and max
random_number() {
  min=$1
  max=$2
  echo $(( $min + RANDOM % ($max - $min + 1) ))
}

# Function to create a test container with a random name
create_test_container() {
  container_name="test-container-$(random_number 1 1000)"
  echo "Creating container: $container_name"
  docker run -d --name $container_name --label chadburn.enabled=true alpine:latest sleep 3600
  echo $container_name
}

# Function to remove a test container
remove_container() {
  container_name=$1
  echo "Removing container: $container_name"
  docker rm -f $container_name >/dev/null 2>&1
}

# Function to execute a command in a container (some will succeed, some will fail)
execute_command() {
  container_name=$1
  should_fail=$(random_number 1 10)
  
  if [ $should_fail -le 2 ]; then
    # 20% chance of failure
    echo "Executing failing command in $container_name"
    docker exec $container_name sh -c "exit 1"
  else
    # 80% chance of success
    echo "Executing successful command in $container_name"
    docker exec $container_name sh -c "echo 'Task executed at $(date)' && sleep $(random_number 1 5)"
  fi
}

# Function to create a Chadburn job-exec for a container
create_job_exec() {
  container_name=$1
  job_id=$(random_number 1 1000)
  
  # Choose a random schedule type
  schedule_type=$(random_number 1 3)
  
  case $schedule_type in
    1)
      # Every X seconds
      interval=$(random_number 10 60)
      schedule="@every ${interval}s"
      ;;
    2)
      # Cron expression (every X minutes)
      minute=$(random_number 0 59)
      schedule="$minute * * * *"
      ;;
    3)
      # Predefined schedule
      predefined=("@hourly" "@daily" "@every 30s" "@every 1m" "@every 5m")
      schedule=${predefined[$(random_number 0 4)]}
      ;;
  esac
  
  echo "Creating job-exec for $container_name with schedule: $schedule"
  
  # Choose a random command type
  command_type=$(random_number 1 4)
  
  case $command_type in
    1)
      # Simple echo
      command="echo 'Scheduled task running at \$(date)'"
      ;;
    2)
      # File operation
      command="touch /tmp/test-file-\$(date +%s) && ls -la /tmp"
      ;;
    3)
      # CPU-intensive task
      command="dd if=/dev/zero bs=1M count=10 | md5sum"
      ;;
    4)
      # Occasionally failing task
      command="if [ \$(( RANDOM % 5 )) -eq 0 ]; then exit 1; else echo 'Success'; fi"
      ;;
  esac
  
  # Add a label to the container for Chadburn to pick up
  docker label $container_name \
    chadburn.job-exec.test-job-$job_id.schedule="$schedule" \
    chadburn.job-exec.test-job-$job_id.command="$command"
}

# Function to create a Chadburn job-local
create_job_local() {
  job_id=$(random_number 1 1000)
  
  # Choose a random schedule
  interval=$(random_number 15 90)
  schedule="@every ${interval}s"
  
  echo "Creating job-local with schedule: $schedule"
  
  # Create a temporary config file with the job-local definition
  cat >> /tmp/job-local-$job_id.conf << EOF
[job-local "test-local-$job_id"]
schedule = $schedule
command = echo "Local job running at \$(date)" && sleep \$(( RANDOM % 5 ))
EOF

  # Copy the config file to the scheduler container
  docker cp /tmp/job-local-$job_id.conf scheduler:/tmp/
  
  # Append the job to the main config file
  docker exec scheduler sh -c "cat /tmp/job-local-$job_id.conf >> /etc/chadburn.conf"
  
  # Remove the temporary file
  rm /tmp/job-local-$job_id.conf
}

# Function to create a Chadburn job-run
create_job_run() {
  job_id=$(random_number 1 1000)
  
  # Choose a random schedule
  interval=$(random_number 20 120)
  schedule="@every ${interval}s"
  
  echo "Creating job-run with schedule: $schedule"
  
  # Create a temporary config file with the job-run definition
  cat >> /tmp/job-run-$job_id.conf << EOF
[job-run "test-run-$job_id"]
schedule = $schedule
image = alpine:latest
command = echo "Run job executing at \$(date)" && sleep \$(( RANDOM % 10 ))
EOF

  # Copy the config file to the scheduler container
  docker cp /tmp/job-run-$job_id.conf scheduler:/tmp/
  
  # Append the job to the main config file
  docker exec scheduler sh -c "cat /tmp/job-run-$job_id.conf >> /etc/chadburn.conf"
  
  # Remove the temporary file
  rm /tmp/job-run-$job_id.conf
}

# Main loop
echo "Starting main loop - press Ctrl+C to stop"
echo "----------------------------------------"

# Keep track of created containers
containers=()

try_count=0
max_tries=100

# Create initial set of containers
for i in {1..3}; do
  new_container=$(create_test_container)
  containers+=($new_container)
  create_job_exec $new_container
done

# Create some job-local and job-run entries
create_job_local
create_job_local
create_job_run
create_job_run

while [ $try_count -lt $max_tries ]; do
  try_count=$((try_count + 1))
  echo "Iteration $try_count of $max_tries"
  
  # Randomly decide what to do
  action=$(random_number 1 10)
  
  if [ $action -le 2 ] && [ ${#containers[@]} -lt 5 ]; then
    # Create a new container (20% chance if less than 5 containers)
    new_container=$(create_test_container)
    containers+=($new_container)
    create_job_exec $new_container
    
  elif [ $action -eq 3 ]; then
    # Create a job-local (10% chance)
    create_job_local
    
  elif [ $action -eq 4 ]; then
    # Create a job-run (10% chance)
    create_job_run
    
  elif [ $action -ge 8 ] && [ ${#containers[@]} -gt 0 ]; then
    # Remove a container (30% chance if there are containers)
    index=$(random_number 0 $((${#containers[@]} - 1)))
    remove_container ${containers[$index]}
    containers=("${containers[@]:0:$index}" "${containers[@]:$((index + 1))}")
    
  elif [ ${#containers[@]} -gt 0 ]; then
    # Execute a command in a random container (30% chance if there are containers)
    index=$(random_number 0 $((${#containers[@]} - 1)))
    execute_command ${containers[$index]}
  fi
  
  # Wait a random amount of time before the next action
  sleep_time=$(random_number 5 15)
  echo "Sleeping for $sleep_time seconds..."
  sleep $sleep_time
  echo ""
done

# Clean up all containers at the end
echo "Cleaning up all test containers..."
for container in "${containers[@]}"; do
  remove_container $container
done

echo "Test data generation complete!" 