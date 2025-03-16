#!/bin/bash

# Verify Chadburn Metrics Script
# This script checks if Chadburn metrics are working correctly

echo "Chadburn Metrics Verification"
echo "============================"
echo

# Check if Chadburn metrics endpoint is accessible
echo "1. Checking if Chadburn metrics endpoint is accessible..."
if curl -s http://localhost:8080/metrics > /dev/null; then
  echo "✅ Chadburn metrics endpoint is accessible"
else
  echo "❌ Chadburn metrics endpoint is not accessible"
  echo "   Make sure Chadburn is running with the --metrics flag and the port is exposed"
  exit 1
fi

# Check if specific metrics are available
echo
echo "2. Checking if Chadburn metrics are available..."
if curl -s http://localhost:8080/metrics | grep -q "chadburn_scheduler_jobs"; then
  echo "✅ Chadburn scheduler metrics are available"
else
  echo "⚠️ Chadburn scheduler metrics are not available"
  echo "   This may be normal if no jobs have been registered yet"
fi

if curl -s http://localhost:8080/metrics | grep -q "chadburn_run_total"; then
  echo "✅ Chadburn job run metrics are available"
else
  echo "⚠️ Chadburn job run metrics are not available"
  echo "   This may be normal if no jobs have run yet"
fi

# Check if Prometheus is accessible
echo
echo "3. Checking if Prometheus is accessible..."
if curl -s http://localhost:9090/-/healthy > /dev/null; then
  echo "✅ Prometheus is accessible"
else
  echo "❌ Prometheus is not accessible"
  echo "   Make sure Prometheus is running and the port is exposed"
  exit 1
fi

# Check if Prometheus is scraping Chadburn metrics
echo
echo "4. Checking if Prometheus is scraping Chadburn metrics..."
if curl -s "http://localhost:9090/api/v1/targets" | grep -q "scheduler:8080"; then
  echo "✅ Prometheus is configured to scrape Chadburn metrics"
else
  echo "❌ Prometheus is not configured to scrape Chadburn metrics"
  echo "   Check your prometheus.yml configuration"
  exit 1
fi

# Check if Grafana is accessible
echo
echo "5. Checking if Grafana is accessible..."
if curl -s http://localhost:3000/api/health > /dev/null; then
  echo "✅ Grafana is accessible"
else
  echo "❌ Grafana is not accessible"
  echo "   Make sure Grafana is running and the port is exposed"
  exit 1
fi

# Check if Prometheus data source is configured in Grafana
echo
echo "6. Checking if Prometheus data source is configured in Grafana..."
if curl -s -H "Authorization: Bearer admin:admin" http://localhost:3000/api/datasources | grep -q "prometheus"; then
  echo "✅ Prometheus data source is configured in Grafana"
else
  echo "❌ Prometheus data source is not configured in Grafana"
  echo "   Run: docker-compose exec datasource-setup /verify-datasource.sh"
  exit 1
fi

# Check if Chadburn dashboard is available in Grafana
echo
echo "7. Checking if Chadburn dashboard is available in Grafana..."
if curl -s -H "Authorization: Bearer admin:admin" http://localhost:3000/api/search?type=dash-db | grep -q "Chadburn"; then
  echo "✅ Chadburn dashboard is available in Grafana"
else
  echo "❌ Chadburn dashboard is not available in Grafana"
  echo "   Check your Grafana provisioning configuration"
  exit 1
fi

# Final summary
echo
echo "============================"
echo "✅ Verification completed successfully!"
echo
echo "You can access the following services:"
echo "- Chadburn Metrics: http://localhost:8080/metrics"
echo "- Prometheus: http://localhost:9090"
echo "- Grafana: http://localhost:3000"
echo
echo "If you don't see any metrics data yet, try running:"
echo "./metrics-tools/add-test-jobs.sh"
echo
echo "This will add test jobs to generate metrics data."
echo "============================" 