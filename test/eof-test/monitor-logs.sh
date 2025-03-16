#!/bin/bash

# This script monitors the Chadburn logs for EOF errors and reconnection attempts

echo "Starting to monitor Chadburn logs..."
echo "Press Ctrl+C to stop monitoring"

# Follow the logs of the Chadburn container
docker logs -f chadburn-test | grep -E "Docker events error|Reconnecting in|Started watching" 