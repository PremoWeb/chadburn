---
layout: default
title: Metrics
nav_order: 4
has_children: true
permalink: /metrics
---

# Chadburn Metrics with Prometheus and Grafana (Experimental)

> **Note**: The metrics functionality in Chadburn is currently experimental. While functional, the API and implementation may change in future releases.
>
> For the most up-to-date documentation on metrics, visit [https://chadburn.dev/metrics](https://chadburn.dev/metrics).

This document explains how to use the preconfigured metrics setup in Chadburn with Prometheus and Grafana for visualization.

## Overview

Chadburn provides a Prometheus-compatible metrics endpoint that exposes various metrics about job executions. This setup includes:

- Chadburn with metrics enabled
- Prometheus for collecting and storing metrics
- Grafana for visualizing metrics with a preconfigured dashboard

## Available Metrics

Chadburn exposes the following metrics:

- `chadburn_scheduler_jobs`: Active job count registered on the scheduler
- `chadburn_scheduler_register_errors_total`: Total number of failed scheduler registrations
- `chadburn_run_total`: Total number of completed job runs (labeled by job name)
- `chadburn_run_errors_total`: Total number of completed job runs that resulted in an error (labeled by job name)
- `chadburn_run_latest_timestamp`: Last time a job run completed (labeled by job name)
- `chadburn_run_duration_seconds`: Duration of all runs (labeled by job name)

## Setup Instructions

### 1. Enable Metrics in Chadburn

To enable metrics in Chadburn, add the `--metrics` flag and specify a listen address:

```bash
chadburn daemon --config=/etc/chadburn.conf --metrics --listen-address=:8080
```

Or in your docker-compose.yml:

```yaml
services:
  scheduler:
    image: premoweb/chadburn:latest
    command: daemon --config=/etc/chadburn.conf --metrics --listen-address=:8080
    ports:
      - "8080:8080"  # Expose the metrics endpoint
```

### 2. Start the Services

The docker-compose.yml file includes all necessary services with preconfigured settings. Start them with:

```bash
docker-compose up -d
```

### 3. Access the Services

Everything is automatically configured, so you can immediately access:

- **Chadburn Metrics**: http://localhost:8080/metrics
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000

Grafana is configured with:
- Anonymous access enabled (no login required)
- The Prometheus data source already configured
- A comprehensive Chadburn dashboard already imported

## Preconfigured Dashboard

The dashboard includes the following panels:

1. **Active Jobs**: Shows the number of jobs registered in the scheduler
2. **Job Runs**: Shows the total number of job runs by job name
3. **Job Errors**: Shows the total number of job errors by job name
4. **Job Duration**: Shows the 95th percentile duration of job runs

## Testing with Sample Data

To generate test data for the dashboard, you can use the included test data generator script:

```bash
./metrics-tools/add-test-jobs.sh
```

This script adds several test jobs to the Chadburn scheduler:

- Regular jobs that run every 10-15 seconds
- A job that occasionally fails (for error metrics)
- A long-running job (for duration metrics)

After running the script, wait a minute for the metrics to appear in Prometheus and Grafana.

You can also verify that your metrics setup is working correctly by running:

```bash
./metrics-tools/verify-metrics.sh
```

This script checks if all components (Chadburn, Prometheus, and Grafana) are properly configured and accessible.

For more information about the testing tools, see the [Test Data Generator](./test-data-generator) documentation.

## Docker Compose Configuration

The complete docker-compose.yml configuration includes:

```yaml
services:
  scheduler:
    build: .
    container_name: scheduler
    restart: unless-stopped
    command: daemon --config=/etc/chadburn.conf --metrics --listen-address=:8080
    ports:
      - "8080:8080"  # Expose the metrics endpoint
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
      - "./test.conf:/etc/chadburn.conf"
      - "./logs:/var/log/chadburn"
    user: "root"
    environment:
      - DOCKER_GID=${DOCKER_GID:-999}
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/metrics"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 5s

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"
    depends_on:
      scheduler:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:9090/-/healthy"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 5s

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_SECURITY_ADMIN_USER=admin
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_INSTALL_PLUGINS=grafana-clock-panel,grafana-simple-json-datasource
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Viewer
      - GF_AUTH_BASIC_ENABLED=false
      - GF_PATHS_PROVISIONING=/etc/grafana/provisioning
    volumes:
      - grafana-storage:/var/lib/grafana
      - ./grafana/provisioning/:/etc/grafana/provisioning/
      - ./grafana/dashboards:/var/lib/grafana/dashboards
    depends_on:
      prometheus:
        condition: service_healthy
    restart: unless-stopped
    
  datasource-setup:
    image: curlimages/curl:latest
    container_name: datasource-setup
    volumes:
      - ./grafana/verify-datasource.sh:/verify-datasource.sh
    entrypoint: ["/bin/sh", "-c"]
    command: ["sleep 15 && /verify-datasource.sh"]
    depends_on:
      grafana:
        condition: service_started
      prometheus:
        condition: service_healthy
```

## Troubleshooting

### No Metrics Showing in Prometheus

1. Check if Chadburn is running with metrics enabled:
   ```bash
   docker-compose logs scheduler
   ```
   Look for a line indicating that metrics are enabled.

2. Check if Prometheus can reach Chadburn:
   ```bash
   docker-compose exec prometheus wget -qO- scheduler:8080/metrics
   ```
   This should display the raw metrics.

3. Check Prometheus targets in the Prometheus UI (Status > Targets) to see if the Chadburn target is up.

### No Data in Grafana

1. Make sure some jobs have run to generate metrics data
2. Check if there are any errors in the Grafana UI
3. Verify that Prometheus is collecting data by checking the Prometheus UI

### Prometheus Data Source Not Working in Grafana

If you see "No data" in the Grafana dashboard panels:

1. Check if the Prometheus data source is working:
   - Go to Grafana (http://localhost:3000)
   - Navigate to Configuration > Data Sources
   - Click on the Prometheus data source
   - Click "Test" at the bottom of the page

2. If the test fails, try restarting the services:
   ```bash
   docker-compose restart prometheus grafana
   ```

3. If issues persist, you can manually fix the data source:
   ```bash
   docker-compose exec datasource-setup /verify-datasource.sh
   ```
   This script will automatically check and fix the Prometheus data source.

4. As a last resort, you can manually configure the data source:
   - Go to Grafana (http://localhost:3000)
   - Navigate to Configuration > Data Sources
   - Click "Add data source"
   - Select "Prometheus"
   - Set the URL to "http://prometheus:9090"
   - Set the UID to "PBFA97CFB590B2093" (important for the dashboard to work)
   - Click "Save & Test"

## Security Considerations

This setup is configured for ease of use with:
- Anonymous access to Grafana (no login required)
- No authentication for the metrics endpoint

In a production environment, consider:

1. Disabling anonymous access to Grafana by setting `GF_AUTH_ANONYMOUS_ENABLED=false` in docker-compose.yml
2. Using a reverse proxy with authentication
3. Restricting network access to the metrics endpoint

## Customizing the Setup

If you want to modify the setup:

1. **Grafana Dashboard**: Edit the dashboard directly in the Grafana UI
2. **Prometheus Configuration**: Modify the prometheus.yml file
3. **Grafana Settings**: Adjust the environment variables in docker-compose.yml

## How It Works

This setup uses Grafana's provisioning feature to automatically:
1. Configure the Prometheus data source
2. Import the Chadburn dashboard
3. Set up anonymous access for easy viewing

All configuration files are stored in the `grafana/` directory.

## Automatic Verification

The setup includes a helper container (`datasource-setup`) that automatically verifies and fixes the Prometheus data source connection. This ensures that the dashboard works correctly without manual intervention.

## Future Improvements

As this is an experimental feature, we plan to make the following improvements:

1. Additional metrics for more detailed job monitoring
2. Better labeling for more granular filtering
3. Integration with other monitoring systems
4. Improved dashboard with more visualization options
5. Support for alerting based on job failures or performance issues

Feedback and contributions to the metrics functionality are welcome! 