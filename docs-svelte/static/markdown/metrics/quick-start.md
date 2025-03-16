---
layout: default
title: Quick Start Guide
parent: Metrics
nav_order: 1
---

# Chadburn Metrics Quick Start (Experimental)

> **Note**: The metrics functionality in Chadburn is currently experimental. While functional, the API and implementation may change in future releases.
>
> For the most up-to-date documentation on metrics, visit [https://chadburn.dev/metrics](https://chadburn.dev/metrics).

This is a preconfigured metrics setup for Chadburn using Prometheus and Grafana. Everything is automatically configured for immediate use.

## Quick Start

1. Use the setup script to set up the Chadburn metrics environment:
   ```bash
   ./metrics-tools/setup-metrics.sh
   ```

   This will start all services, verify the setup, and optionally add test jobs.

2. Alternatively, you can start the services manually:
   ```bash
   docker-compose up -d
   ```

3. Access the dashboards:
   - **Grafana Dashboard**: http://localhost:3000
   - **Prometheus**: http://localhost:9090
   - **Raw Metrics**: http://localhost:8080/metrics

## What's Included

- **Chadburn** with metrics enabled
- **Prometheus** for collecting metrics
- **Grafana** with:
  - Preconfigured data source
  - Comprehensive dashboard
  - Anonymous access (no login required)
- **Automatic verification** that ensures everything works correctly
- **Testing tools** in the `metrics-tools/` directory:
  - `setup-metrics.sh`: Sets up the Chadburn metrics environment
  - `verify-metrics.sh`: Verifies that the setup is working correctly
  - `add-test-jobs.sh`: Adds test jobs to generate metrics data

## Dashboard Features

The preconfigured dashboard shows:
- Active job count
- Job runs by name
- Job errors by name
- Job execution duration

## Generate Test Data

To generate test data for the dashboard, run:

```bash
./metrics-tools/add-test-jobs.sh
```

This script adds several test jobs to the Chadburn scheduler that will generate metrics data. Wait a minute for the metrics to appear in Prometheus and Grafana.

## Verify Your Setup

To verify that your metrics setup is working correctly, run:

```bash
./metrics-tools/verify-metrics.sh
```

This script checks if all components (Chadburn, Prometheus, and Grafana) are properly configured and accessible.

## Troubleshooting

If you encounter any issues with the dashboard not showing data:

1. Wait a minute for the automatic verification to complete
2. If issues persist, run: `docker-compose exec datasource-setup /verify-datasource.sh`
3. Make sure some jobs have run to generate metrics data

## More Information

For detailed information about the metrics setup, see the [main metrics documentation](./index).

For information about the testing tools, see the [Test Data Generator](./test-data-generator) documentation. 