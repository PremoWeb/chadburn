#!/bin/bash

# Script to add test jobs directly to the scheduler container

echo "Adding test jobs to Chadburn scheduler..."

# Create a temporary config file with test jobs
cat > test-jobs.conf << EOF
[job-local "test-job-1"]
schedule = @every 10s
command = echo "Test job 1 running at \$(date)" && sleep 2

[job-local "test-job-2"]
schedule = @every 15s
command = echo "Test job 2 running at \$(date)" && sleep 3

[job-local "failing-job-1"]
schedule = @every 20s
command = if [ \$(( RANDOM % 3 )) -eq 0 ]; then echo "Job failed" && exit 1; else echo "Job succeeded" && sleep 1; fi

[job-local "long-job-1"]
schedule = @every 30s
command = echo "Starting long job" && sleep 8 && echo "Finished long job"

[job-local "test-job-3"]
schedule = @every 12s
command = echo "Test job 3 running at \$(date)" && sleep 1
EOF

# Copy the config file to the scheduler container
docker cp test-jobs.conf scheduler:/tmp/

# Append the jobs to the main config file
docker exec scheduler sh -c "cat /tmp/test-jobs.conf >> /etc/chadburn.conf"

# Restart the scheduler to pick up the new jobs
docker restart scheduler

echo "Test jobs added. Wait a minute for metrics to appear in Prometheus and Grafana."
echo "You can check the metrics at:"
echo "- Grafana: http://localhost:3000"
echo "- Prometheus: http://localhost:9090"
echo "- Raw metrics: http://localhost:8080/metrics" 