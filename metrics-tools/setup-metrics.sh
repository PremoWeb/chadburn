#!/bin/bash

# Chadburn Metrics Setup Script
# This script helps set up the Chadburn metrics environment

echo "Chadburn Metrics Setup"
echo "====================="
echo

# Check if Docker is installed
echo "Checking if Docker is installed..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi
echo "✅ Docker is installed"

# Check if Docker Compose is installed
echo "Checking if Docker Compose is installed..."
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi
echo "✅ Docker Compose is installed"

# Start the services
echo
echo "Starting Chadburn metrics services..."
docker-compose up -d
echo "✅ Services started"

# Wait for services to be ready
echo
echo "Waiting for services to be ready..."
sleep 15
echo "✅ Services should be ready now"

# Verify the setup
echo
echo "Verifying the setup..."
./verify-metrics.sh

# Add test jobs
echo
echo "Would you like to add test jobs to generate metrics data? (y/n)"
read -r response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    ./add-test-jobs.sh
    echo "✅ Test jobs added"
    echo "Wait a minute for metrics to appear in Prometheus and Grafana"
else
    echo "Skipping test jobs"
    echo "You can add test jobs later by running: ./add-test-jobs.sh"
fi

echo
echo "====================="
echo "Setup completed!"
echo
echo "You can access the following services:"
echo "- Chadburn Metrics: http://localhost:8080/metrics"
echo "- Prometheus: http://localhost:9090"
echo "- Grafana: http://localhost:3000"
echo
echo "For more information, see:"
echo "- README-METRICS.md for a quick start guide"
echo "- METRICS.md for detailed documentation"
echo "=====================" 