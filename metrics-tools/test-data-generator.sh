#!/bin/bash

# Test Data Generator for Chadburn Metrics Dashboard
# This script randomly starts/stops containers and executes tasks to generate metrics

echo "Starting Chadburn Test Data Generator"
echo "======================================"

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

# Function to create a Chadburn job for a container
create_chadburn_job() {
  container_name=$1
  job_id=$(random_number 1 1000)
  interval=$(random_number 10 60)
  
  echo "Creating Chadburn job for $container_name (runs every $interval seconds)"
  
  # Add a label to the container for Chadburn to pick up
  docker label $container_name \
    chadburn.job-exec.test-job-$job_id.schedule="@every ${interval}s" \
    chadburn.job-exec.test-job-$job_id.command="echo 'Scheduled task running at \$(date)' && sleep \$(( RANDOM % 10 ))"
}

# Main loop
echo "Starting main loop - press Ctrl+C to stop"
echo "----------------------------------------"

# Keep track of created containers
containers=()

try_count=0
max_tries=100

while [ $try_count -lt $max_tries ]; do
  try_count=$((try_count + 1))
  echo "Iteration $try_count of $max_tries"
  
  # Randomly decide what to do
  action=$(random_number 1 10)
  
  if [ $action -le 3 ] && [ ${#containers[@]} -lt 5 ]; then
    # Create a new container (30% chance if less than 5 containers)
    new_container=$(create_test_container)
    containers+=($new_container)
    create_chadburn_job $new_container
    
  elif [ $action -ge 8 ] && [ ${#containers[@]} -gt 0 ]; then
    # Remove a container (30% chance if there are containers)
    index=$(random_number 0 $((${#containers[@]} - 1)))
    remove_container ${containers[$index]}
    containers=("${containers[@]:0:$index}" "${containers[@]:$((index + 1))}")
    
  elif [ ${#containers[@]} -gt 0 ]; then
    # Execute a command in a random container (40% chance if there are containers)
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