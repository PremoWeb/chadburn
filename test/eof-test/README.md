# Docker Events EOF Error Test

This test environment is designed to verify the fix for issue #115, which addresses the "Error watching events: unexpected EOF" issue that can lead to memory problems and container restarts.

## Test Setup

The test environment consists of:

1. A Chadburn container built with the EOF error fix
2. A test container with the `chadburn.enabled=true` label
3. A script to simulate EOF errors by restarting the Docker daemon
4. A script to monitor Chadburn logs

## Running the Test

### Step 1: Build and Start the Test Environment

```bash
cd /home/nick/Projects/chadburn/test/eof-test
docker-compose up -d
```

### Step 2: Monitor Chadburn Logs

In one terminal, run:

```bash
./monitor-logs.sh
```

### Step 3: Simulate EOF Errors

In another terminal, run:

```bash
./simulate-eof.sh
```

This script requires sudo privileges to restart the Docker daemon.

### Step 4: Observe the Results

Watch the logs to see how Chadburn handles the EOF errors. You should see:

1. EOF errors being reported
2. Log messages about reconnecting with increasing delays
3. Successful reconnections after the Docker daemon restarts
4. No memory issues or container restarts

## Expected Behavior

With the fix applied, Chadburn should:

1. Detect EOF errors
2. Apply exponential backoff when reconnecting
3. Reset the backoff timer when events are successfully received
4. Continue operating normally without memory issues

## Cleanup

When you're done testing, stop the test environment:

```bash
docker-compose down
``` 