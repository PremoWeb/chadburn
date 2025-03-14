# Chadburn: A Modern Job Scheduler

[![GitHub release](https://img.shields.io/github/v/release/PremoWeb/Chadburn)](https://github.com/PremoWeb/Chadburn/releases)
[![Testing Status](https://github.com/PremoWeb/Chadburn/workflows/Testing%20Status/badge.svg)](https://github.com/PremoWeb/Chadburn/actions?query=workflow%3A%22Testing+Status%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/premoweb/chadburn)](https://hub.docker.com/r/premoweb/chadburn)
[![Docker Image Size](https://img.shields.io/docker/image-size/premoweb/chadburn/latest)](https://hub.docker.com/r/premoweb/chadburn)
[![License](https://img.shields.io/github/license/PremoWeb/Chadburn)](https://github.com/PremoWeb/Chadburn/blob/main/LICENSE)

**Chadburn** is a lightweight job scheduler designed for __Docker__ environments, developed in Go. It serves as a contemporary replacement for the traditional [cron](https://en.wikipedia.org/wiki/Cron).

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

### Key Features

Chadburn's primary feature is its ability to execute commands directly within Docker containers. Utilizing Docker's API, Chadburn mimics the behavior of [`exec`](https://docs.docker.com/reference/commandline/exec/), enabling commands to run inside active containers. Additionally, it allows for command execution in new containers, which are destroyed after use.

---

## Configuration

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

For detailed parameters, refer to the [Jobs reference documentation](docs/jobs.md).

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

You can also add `job-run` labels directly to containers you want to start periodically. Chadburn will automatically detect these containers and schedule them to start according to the specified schedule:

```bash
docker run -d --name my-periodic-container \
    --label chadburn.job-run.schedule="@daily" \
    my-image
```

This will create a job that starts the `my-periodic-container` container once a day. The container will be started using the same configuration it was created with.

### Logging

Chadburn offers three logging drivers that can be configured in the `[global]` section:

- `mail` to send notifications via email.
- `save` to save structured execution reports in a specified directory.
- `slack` to send messages through a Slack webhook.

All logs include timestamps in the format `YYYY-MM-DD HH:MM:SS.mmm`. You can also use Docker's built-in timestamp feature with `docker logs --timestamps chadburn` if you prefer Docker's timestamp format.

#### Logging Options

- `smtp-host`,