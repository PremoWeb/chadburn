#!/bin/bash

# Wait for Grafana to be ready
echo "Waiting for Grafana to be ready..."
until $(curl --output /dev/null --silent --head --fail http://grafana:3000/api/health); do
  printf '.'
  sleep 5
done
echo "Grafana is up!"

# Wait for Prometheus to be ready
echo "Waiting for Prometheus to be ready..."
until $(curl --output /dev/null --silent --head --fail http://prometheus:9090/-/healthy); do
  printf '.'
  sleep 5
done
echo "Prometheus is up!"

# Check if the Prometheus data source exists and is working
echo "Checking Prometheus data source..."
DATASOURCE_STATUS=$(curl -s -H "Content-Type: application/json" -X GET http://grafana:3000/api/datasources/uid/PBFA97CFB590B2093/health)

if [[ $DATASOURCE_STATUS == *"success"* ]]; then
  echo "Prometheus data source is working correctly!"
else
  echo "Prometheus data source is not working. Attempting to fix..."
  
  # Delete the data source if it exists but is not working
  curl -s -X DELETE http://grafana:3000/api/datasources/uid/PBFA97CFB590B2093
  
  # Create the data source with the correct settings
  curl -s -X POST -H "Content-Type: application/json" -d '{
    "name": "Prometheus",
    "type": "prometheus",
    "url": "http://prometheus:9090",
    "access": "proxy",
    "isDefault": true,
    "uid": "PBFA97CFB590B2093"
  }' http://grafana:3000/api/datasources
  
  echo "Data source has been recreated."
fi

echo "Setup complete!" 