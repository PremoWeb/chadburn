# Chadburn Metrics Tools

This directory contains tools for testing and verifying the Chadburn metrics functionality.

## Available Tools

### setup-metrics.sh

This script helps set up the Chadburn metrics environment. It:

- Checks if Docker and Docker Compose are installed
- Starts the Chadburn metrics services
- Verifies that the setup is working correctly
- Optionally adds test jobs to generate metrics data

Usage:
```bash
./setup-metrics.sh
```

### add-test-jobs.sh

This script adds test jobs to the Chadburn scheduler for generating metrics data. It creates:

- Regular jobs that run every 10-15 seconds
- A job that occasionally fails (for error metrics)
- A long-running job (for duration metrics)

Usage:
```bash
./add-test-jobs.sh
```

### verify-metrics.sh

This script verifies that your metrics setup is working correctly. It checks:

- If the Chadburn metrics endpoint is accessible
- If specific metrics are available
- If Prometheus is accessible and configured to scrape Chadburn metrics
- If Grafana is accessible and configured with the Prometheus data source
- If the Chadburn dashboard is available in Grafana

Usage:
```bash
./verify-metrics.sh
```

### simple-test-generator.sh

This script is used by the `test-data-generator` service in the docker-compose.yml file. It adds test jobs directly to the Chadburn configuration file inside the container.

### test-data-generator-advanced.sh

This is an advanced version of the test data generator that provides more options and flexibility for generating test data.

## Documentation

For comprehensive documentation on Chadburn metrics, please see:

- Our website at [https://chadburn.dev/metrics](https://chadburn.dev/metrics) for the latest documentation 