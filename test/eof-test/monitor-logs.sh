#!/bin/bash

# This script monitors the Chadburn logs for EOF errors and reconnection attempts

echo "Starting log monitoring..."
echo "Press Ctrl+C to stop monitoring."

docker logs -f chadburn-test | grep -E "EOF|Docker events|error|reconnect" 