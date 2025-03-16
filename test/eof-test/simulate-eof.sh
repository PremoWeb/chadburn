#!/bin/bash

# This script simulates EOF errors by restarting the Docker daemon
# Note: This requires sudo privileges

echo "Starting EOF error simulation..."
echo "This will restart the Docker daemon to simulate EOF errors."
echo "Press Ctrl+C to stop the simulation."

while true; do
  echo "Simulating EOF error by restarting Docker daemon..."
  sudo systemctl restart docker
  echo "Docker daemon restarted. Waiting for 30 seconds..."
  sleep 30
done 