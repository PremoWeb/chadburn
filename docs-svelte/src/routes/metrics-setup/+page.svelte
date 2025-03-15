<script lang="ts">
  import { base } from '$app/paths';
</script>

<svelte:head>
  <title>Chadburn - Metrics Setup</title>
</svelte:head>

<div class="content">
  <h1>Setting Up Chadburn's Metrics Endpoint</h1>
  
  <p>
    Chadburn includes a built-in HTTP server that exposes Prometheus-compatible metrics. 
    This allows you to monitor job executions, errors, and performance using tools like 
    Prometheus and Grafana.
  </p>

  <h2>Enabling the Metrics Endpoint</h2>
  
  <p>
    To enable the metrics endpoint, use the <code>--metrics</code> flag when starting Chadburn:
  </p>

  <pre><code>chadburn daemon --config=/etc/chadburn.conf --metrics</code></pre>

  <p>
    By default, the metrics endpoint listens on port 8080. You can change this using the 
    <code>--listen-address</code> flag:
  </p>

  <pre><code>chadburn daemon --config=/etc/chadburn.conf --metrics --listen-address=:9100</code></pre>

  <h2>Available Metrics</h2>

  <p>The following metrics are available:</p>

  <ul>
    <li><code>chadburn_scheduler_jobs</code>: Active job count registered on the scheduler</li>
    <li><code>chadburn_scheduler_register_errors_total</code>: Total number of failed scheduler registrations</li>
    <li><code>chadburn_run_total</code>: Total number of completed job runs (labeled by job name)</li>
    <li><code>chadburn_run_errors_total</code>: Total number of completed job runs that resulted in an error (labeled by job name)</li>
    <li><code>chadburn_run_latest_timestamp</code>: Last time a job run completed (labeled by job name)</li>
    <li><code>chadburn_run_duration_seconds</code>: Duration of all runs (labeled by job name)</li>
  </ul>

  <h2>Docker Compose Setup</h2>

  <p>
    Here's a complete setup using Docker Compose that includes Chadburn with metrics enabled, 
    Prometheus for collecting metrics, and Grafana for visualization:
  </p>

  <pre><code>{`version: '3'

services:
  chadburn:
    image: premoweb/chadburn:latest
    container_name: chadburn
    volumes:
      - ./example-config.ini:/etc/chadburn.conf
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
      - /var/log/chadburn:/var/log/chadburn
    ports:
      - "8080:8080"  # Expose the metrics endpoint
    command: daemon --config=/etc/chadburn.conf --metrics --listen-address=:8080
    restart: unless-stopped

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"
    depends_on:
      - chadburn
    restart: unless-stopped

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana
    depends_on:
      - prometheus
    restart: unless-stopped

volumes:
  grafana-storage:`}</code></pre>

  <h3>Prometheus Configuration</h3>

  <p>Create a <code>prometheus.yml</code> file with the following content:</p>

  <pre><code>{`global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'chadburn'
    static_configs:
      - targets: ['chadburn:8080']`}</code></pre>

  <h2>Accessing the Metrics</h2>

  <p>
    Once Chadburn is running with metrics enabled, you can access the metrics directly at:
  </p>

  <pre><code>http://localhost:8080/metrics</code></pre>

  <p>
    This will display all the Prometheus-compatible metrics in plain text format.
  </p>

  <h2>Example Configuration Files</h2>

  <p>
    You can download the example configuration files here:
  </p>

  <ul>
    <li><a href="{base}/data/docker-compose.yml" download>docker-compose.yml</a></li>
    <li><a href="{base}/data/prometheus.yml" download>prometheus.yml</a></li>
    <li><a href="{base}/data/example-config.ini" download>example-config.ini</a></li>
  </ul>

  <h2>Security Considerations</h2>

  <p>
    The metrics endpoint does not include authentication by default. In production environments, 
    you should:
  </p>

  <ul>
    <li>Use a reverse proxy with authentication in front of the metrics endpoint</li>
    <li>Use network policies to restrict access to the metrics endpoint</li>
    <li>Consider binding to localhost only and using SSH tunneling for remote access</li>
  </ul>
</div>

<style>
  .content {
    max-width: 800px;
    margin: 0 auto;
  }

  pre {
    background-color: #f5f5f5;
    padding: 1rem;
    border-radius: 4px;
    overflow-x: auto;
  }

  code {
    font-family: monospace;
  }
</style> 