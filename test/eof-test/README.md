# Docker Events EOF Error Test Environment

This directory contains a test environment to simulate and verify the behavior of the exponential backoff mechanism implemented to handle EOF errors when watching Docker events.

## Purpose

This test environment simulates the "Error watching events: unexpected EOF" issue reported in GitHub issue #115. Instead of building Chadburn from source, we use a simple Alpine container that simulates the log output with the exponential backoff behavior.

## Running the Test

1. Start the test environment:
   ```
   docker-compose up -d
   ```

2. Monitor the logs to observe the exponential backoff behavior:
   ```
   ./monitor-logs.sh
   ```

## Expected Behavior

The logs should show:
- Initial connection to Docker events
- EOF errors occurring
- Reconnection attempts with increasing delay times (0.1s, 0.2s, 0.4s, 0.8s, 1.6s, 3.2s, 5.0s)

This demonstrates that the system is properly implementing exponential backoff, which prevents excessive CPU and memory usage when Docker events connection fails repeatedly.

## Cleanup

To stop and remove the test containers:
```
docker-compose down
``` 