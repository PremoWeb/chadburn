---
layout: default
title: Test Data Generator
parent: Metrics
nav_order: 2
---

# Chadburn Metrics Test Data Generator

This document explains the test data generator tools available for Chadburn metrics.

## Overview

The test data generator tools are designed to help you generate test data for the Chadburn metrics dashboard. These tools create various types of jobs that run at different intervals, including jobs that occasionally fail and long-running jobs, to provide a comprehensive view of the metrics capabilities.

## Available Tools

All testing tools are located in the `metrics-tools/` directory:

### setup-metrics.sh

This script helps set up the Chadburn metrics environment. It:

- Checks if Docker and Docker Compose are installed
- Starts the Chadburn metrics services
- Verifies that the setup is working correctly
- Optionally adds test jobs to generate metrics data

Usage:
```bash
./metrics-tools/setup-metrics.sh
```

### add-test-jobs.sh

This script adds test jobs to the Chadburn scheduler for generating metrics data. It creates:

- Regular jobs that run every 10-15 seconds
- A job that occasionally fails (for error metrics)
- A long-running job (for duration metrics)

Usage:
```bash
./metrics-tools/add-test-jobs.sh
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
./metrics-tools/verify-metrics.sh
```

### simple-test-generator.sh

This script is used by the `test-data-generator` service in the docker-compose.yml file. It adds test jobs directly to the Chadburn configuration file inside the container.

### test-data-generator-advanced.sh

This is an advanced version of the test data generator that provides more options and flexibility for generating test data.

## How to Use

1. Run the setup script to set up the Chadburn metrics environment:
   ```bash
   ./metrics-tools/setup-metrics.sh
   ```

   This will start all services, verify the setup, and optionally add test jobs.

2. Alternatively, you can perform these steps manually:

   a. Start the Chadburn metrics setup:
      ```bash
      docker-compose up -d
      ```

   b. Verify that the setup is working correctly:
      ```bash
      ./metrics-tools/verify-metrics.sh
      ```

   c. Generate test data:
      ```bash
      ./metrics-tools/add-test-jobs.sh
      ```

3. Access the dashboards:
   - Grafana: http://localhost:3000
   - Prometheus: http://localhost:9090
   - Raw Metrics: http://localhost:8080/metrics

## Job Types

The test data generator creates the following types of jobs:

1. **Regular Jobs**: These jobs run at regular intervals (every 10-15 seconds) and complete successfully.
2. **Failing Jobs**: These jobs occasionally fail (with a configurable failure rate) to generate error metrics.
3. **Long-Running Jobs**: These jobs take longer to complete (5-15 seconds) to demonstrate duration metrics.

## Configuration

The test jobs are configured with the following parameters:

- **Schedule**: How often the job runs (e.g., `@every 10s`)
- **Command**: The command to execute
- **Failure Rate**: For failing jobs, the probability of failure (e.g., 20%)
- **Duration**: For long-running jobs, how long the job takes to complete

## Customization

You can customize the test data generator by editing the scripts in the `metrics-tools/` directory. For example, you can:

- Change the job schedules
- Modify the failure rates
- Add more job types
- Adjust the job durations

## More Information

For more information about Chadburn metrics, see the [main metrics documentation](./index). 