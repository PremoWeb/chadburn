# Chadburn: A Modern Job Scheduler

[![GitHub release](https://img.shields.io/github/v/release/PremoWeb/Chadburn)](https://github.com/PremoWeb/Chadburn/releases)
[![Testing Status](https://github.com/PremoWeb/Chadburn/workflows/Testing%20Status/badge.svg)](https://github.com/PremoWeb/Chadburn/actions?query=workflow%3A%22Testing+Status%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/premoweb/chadburn)](https://hub.docker.com/r/premoweb/chadburn)
[![Docker Image Size](https://img.shields.io/docker/image-size/premoweb/chadburn/latest)](https://hub.docker.com/r/premoweb/chadburn)
[![License](https://img.shields.io/github/license/PremoWeb/Chadburn)](https://github.com/PremoWeb/Chadburn/blob/main/LICENSE)

> ### 📣 **NEW!** Visit our dedicated website at [https://chadburn.dev](https://chadburn.dev) for comprehensive documentation, tutorials, and examples. Our documentation is now powered by SvelteKit and hosted on Cloudflare Pages for improved performance and reliability.

**Chadburn** is a lightweight job scheduler designed for __Docker__ environments, developed in Go. It serves as a contemporary replacement for the traditional [cron](https://en.wikipedia.org/wiki/Cron). Chadburn uses semantic versioning and automatic releases based on commit messages. Commit messages are validated using Git hooks to ensure they follow the Conventional Commits format.

---

### Special Note

Chadburn is a project built upon the ongoing development from Ofelia, a fork initiated by @rdelcorro. This project was created to address specific needs, including:

- Automatic task updates when Docker containers are started, stopped, restarted, or modified.
- Elimination of the need for a dummy task in the Chadburn container.
- Concurrent support for both INI files and Docker labels, allowing configurations to be merged seamlessly.
- The ability to recognize new tasks or remove existing ones without needing a restart.

---

### Why Choose Chadburn?

Since the release of [`cron`](https://en.wikipedia.org/wiki/Cron) by AT&T Bell Laboratories in March 1975, much has changed in the computing landscape, especially with the rise of Docker. While **Vixie's cron** remains functional, it lacks extensibility and can be challenging to debug when issues arise.

Various solutions exist, including containerized cron implementations and command wrappers, but these often complicate straightforward tasks.

---

### Important Update: Migration to Official Docker Client

Chadburn has removed all references to the legacy polyfill dependency `fsouza/go-dockerclient` and migrated to the official Docker client library. This change brings improved compatibility with the latest Docker features and better long-term maintainability.

We would like to express our sincere gratitude to Francisco Souza and all contributors to the `fsouza/go-dockerclient` library. Their exceptional work provided a robust foundation for Chadburn and many other Docker-related projects over the years. The library's reliability and comprehensive API coverage were instrumental in Chadburn's early development and success.

---

### Key Features

Chadburn's primary feature is its ability to execute commands directly within Docker containers. Utilizing Docker's API, Chadburn mimics the behavior of [`exec`](https://docs.docker.com/reference/commandline/exec/), enabling commands to run inside active containers. Additionally, it allows for command execution in new containers, which are destroyed after use.

Chadburn also supports variable substitution in job commands, allowing you to reference container information dynamically using syntax like `{{.Container.Name}}` and `{{.Container.ID}}`. This makes it easier to create reusable job configurations that can interact with containers without hardcoding their names or IDs.

---

## Configuration

For the most comprehensive and up-to-date documentation, please visit our dedicated website at [https://chadburn.dev](https://chadburn.dev).

A comprehensive wiki is underway to detail Chadburn's usage. Caprover users will soon have access to a One Click App for deploying and managing scheduled jobs via Service Label Overrides.

For others, here's a quick guide to get started with Chadburn:

### Job Scheduling

Chadburn uses a scheduling format consistent with the Go implementation of `cron`. Examples include `@every 10s` or `0 0 1 * * *` (which runs every night at 1 AM).

**Note**: The scheduling format previously included seconds; however, this has been updated in the latest version of Chadburn. Significant development is planned to resolve various issues reported with both Ofelia and Chadburn.

You can configure four types of jobs:

- `job-exec`: Executes a command inside a running container.
- `job-run`: Runs a command in a new container using a specified image.
- `job-local`: Executes a command on the host running Chadburn.
- `job-service-run`: Runs a command inside a new "run-once" service for swarm environments.

For detailed parameters, refer to the [Jobs reference documentation](https://chadburn.dev/jobs).

#### INI Configuration

To run Chadburn with an INI file, use the command:

```bash
chadburn daemon --config=/path/to/config.ini
```

Here's a sample INI configuration:

```ini
[job-exec "job-executed-on-running-container"]
schedule = @hourly
container = my-container
command = touch /tmp/example

[job-run "job-executed-on-new-container"]
schedule = @hourly
image = ubuntu:latest
command = touch /tmp/example

[job-local "job-executed-on-current-host"]
schedule = @hourly
command = touch /tmp/example

[job-service-run "service-executed-on-new-container"]
schedule = 0,20,40 * * * *
image = ubuntu
network = swarm_network
command = touch /tmp/example
```

#### Docker Label Configurations

For Docker label configurations, Chadburn needs access to the Docker socket:

```bash
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    premoweb/chadburn:latest daemon
```

> **Note**: If you encounter permission issues with the Docker socket, refer to our [Docker Socket Permissions Guide](https://chadburn.dev/docker-socket-permissions) for solutions.

The labels format is: `chadburn.<JOB_TYPE>.<JOB_NAME>.<JOB_PARAMETER>=<PARAMETER_VALUE>`. This configuration method supports all capabilities provided by INI files.

To execute `job-exec`, the target container must have the label `chadburn.enabled=true`.

For example, to run the `uname -a` command in an existing container called `my_nginx`, start `my_nginx` with the following configurations:

```bash
docker run -it --rm \
    --label chadburn.enabled=true \
    --label chadburn.job-exec.test-exec-job.schedule="@every 5s" \
    --label chadburn.job-exec.test-exec-job.command="uname -a" \
    nginx
```

Alternatively, you can use Docker Compose:

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    depends_on:
      - nginx
    command: daemon
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro

  nginx:
    image: nginx
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.datecron.schedule: "@every 5s"
      chadburn.job-exec.datecron.command: "uname -a"
```

#### Dynamic Docker Configuration

Chadburn can be run in its own container or directly on the host. It will automatically detect any containers that start, stop, or change, utilizing the labeled containers for dynamic task management.

#### Hybrid Configuration (INI + Docker)

You can combine INI files and Docker labels to manage configurations. Use INI files for global settings or tasks that cannot be defined solely through labels, while Docker labels can be employed for dynamically managed tasks.

```ini
[global]
slack-webhook = https://myhook.com/auth

[job-run "job-executed-on-new-container"]
schedule = @hourly
image = ubuntu:latest
command = touch /tmp/example
```

Use Docker labels for dynamic jobs:

```bash
docker run -it --rm \
    --label chadburn.enabled=true \
    --label chadburn.job-exec.test-exec-job.schedule="@every 5s" \
    --label chadburn.job-exec.test-exec-job.command="uname -a" \
    nginx
```

#### Running Containers with Labels

You can add `job-run` labels directly to containers you want to start periodically. Chadburn will automatically detect these containers and schedule them to start according to the specified schedule:

```bash
docker run -d --name my-periodic-container \
    --label chadburn.job-run.schedule="@daily" \
    my-image
```

This will create a job that starts the `my-periodic-container` container once a day. The container will be started using the same configuration it was created with.

> **Important Note about job-run with Docker Compose**: When using `job-run` with Docker Compose, there's a key difference from `job-exec`. The `job-run` type is designed to **start new containers** or **start existing stopped containers**, not to execute commands in already running containers. 
>
> For example, this configuration will **NOT** run the command every 5 seconds inside the running container:
> ```yaml
> services:
>   alpine:
>     image: alpine
>     labels:
>       chadburn.enabled: "true"
>       chadburn.job-run.datecron.schedule: "@every 5s"
>       chadburn.job-run.datecron.command: "uname -a"
> ```
>
> Instead, if you want to run a command periodically inside a running container, use `job-exec`:
> ```yaml
> services:
>   alpine:
>     image: alpine
>     labels:
>       chadburn.enabled: "true"
>       chadburn.job-exec.datecron.schedule: "@every 5s"
>       chadburn.job-exec.datecron.command: "uname -a"
> ```

### Logging

Chadburn offers three logging drivers that can be configured in the `[global]` section:

- `mail` to send notifications via email.
- `save` to save structured execution reports in a specified directory.
- `slack` to send messages through a Slack webhook.

All logs include timestamps in the format `YYYY-MM-DD HH:MM:SS.mmm`. You can also use Docker's built-in timestamp feature with `docker logs --timestamps chadburn` if you prefer Docker's timestamp format.

#### Logging Options

- `smtp-host`, `smtp-port`, `smtp-user`, `smtp-password`: SMTP server configuration for email notifications.
- `email-to`, `email-from`: Email addresses for notifications.
- `mail-only-on-error`: If set to true, emails are only sent when a job fails.
- `insecure-skip-verify`: If set to true, skips TLS certificate verification for SMTP.
- `slack-webhook`: Webhook URL for Slack notifications.
- `slack-only-on-error`: If set to true, Slack messages are only sent when a job fails.
- `save-folder`: Directory where execution logs should be saved.
- `save-only-on-error`: If set to true, logs are only saved when a job fails.

Example configuration with logging:

```ini
[global]
smtp-host = smtp.example.com
smtp-port = 587
smtp-user = user
smtp-password = password
email-to = alerts@example.com
email-from = chadburn@example.com
mail-only-on-error = true

slack-webhook = https://hooks.slack.com/services/XXX/YYY/ZZZ
slack-only-on-error = true

save-folder = /var/log/chadburn
save-only-on-error = false
```

### Metrics (Experimental)

Chadburn includes experimental support for Prometheus metrics, allowing you to monitor job executions and performance. When enabled, Chadburn exposes a metrics endpoint that can be scraped by Prometheus.

To enable metrics, add the `--metrics` flag and specify a listen address:

```bash
chadburn daemon --config=/etc/chadburn.conf --metrics --listen-address=:8080
```

Available metrics include:
- Job counts
- Job execution totals
- Error counts
- Execution durations

A preconfigured setup with Prometheus and Grafana is included for easy visualization of metrics. Testing and verification tools are available in the `metrics-tools/` directory. For more information, see:
- The [metrics documentation](https://chadburn.dev/metrics) for comprehensive information about Chadburn's metrics capabilities

> **Note**: The metrics functionality is currently experimental and may change in future releases.

## Installation

The simplest way to deploy **Chadburn** is using Docker, as outlined above.

If you prefer not to use the provided Docker image, you can download a pre-built binary from the [releases page](https://github.com/PremoWeb/chadburn/releases). Chadburn provides binaries for multiple platforms:

- **Linux**: amd64, arm64, and armv7
- **macOS**: amd64 and arm64 (Apple Silicon)
- **Windows**: amd64 and arm64

Each release includes SHA256 checksums for verifying the integrity of the downloaded binaries.

### Installing from Binary

1. Download the appropriate binary for your platform from the [releases page](https://github.com/PremoWeb/chadburn/releases).
2. Extract the archive:
   ```bash
   # For Linux/macOS
   tar -xzf chadburn-<platform>-<arch>.tar.gz
   
   # For Windows
   # Extract the .zip file using your preferred tool
   ```
3. Move the binary to a location in your PATH:
   ```bash
   # For Linux/macOS
   sudo mv chadburn-<platform>-<arch> /usr/local/bin/chadburn
   chmod +x /usr/local/bin/chadburn
   
   # For Windows
   # Move the .exe file to a suitable location and add it to your PATH
   ```

### Building from Source

If you prefer to build from source, you'll need Go 1.19 or later:

```bash
git clone https://github.com/PremoWeb/chadburn.git
cd chadburn
go build -o chadburn
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

For contribution guidelines, development setup, and best practices, visit our [contribution page](https://chadburn.dev/contributing) on our website.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Ofelia](https://github.com/mcuadros/ofelia) - The original project that inspired Chadburn.
- [Cron](https://github.com/robfig/cron) - The Go cron library used by Chadburn.
- [Docker](https://www.docker.com/) - The container platform that Chadburn integrates with.
