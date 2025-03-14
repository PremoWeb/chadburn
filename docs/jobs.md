# Jobs reference

- [job-exec](#job-exec)
- [job-run](#job-run)
- [job-local](#job-local)
- [job-service-run](#job-service-run)
- [Variable Substitution](#variable-substitution)

## Job Types

Chadburn supports several job types:

- `job-local`: Executes commands on the host machine.
- `job-exec`: Executes commands inside a running container.
- `job-run`: Creates a new container to execute commands.
- `job-service-run`: Creates a service to execute commands.
- `job-lifecycle`: Executes commands on container lifecycle events (start/stop).

### job-lifecycle

The `job-lifecycle` type allows you to execute commands when a container starts or stops. This is useful for sending notifications, performing cleanup, or triggering other actions based on container lifecycle events.

Unlike other job types, `job-lifecycle` jobs don't run on a schedule. Instead, they run once when the specified container lifecycle event occurs.

#### Configuration

```toml
[job-lifecycle "notify-on-start"]
container = "my-container"
event-type = "start"
command = "echo 'Container my-container started' | mail -s 'Container Started' admin@example.com"

[job-lifecycle "cleanup-on-stop"]
container = "my-container"
event-type = "stop"
command = "echo 'Container my-container stopped' | mail -s 'Container Stopped' admin@example.com"
```

#### Docker Compose Labels

```yaml
services:
  my-service:
    image: nginx
    labels:
      chadburn.enabled: "true"
      chadburn.job-lifecycle.notify-on-start.command: "echo 'Container {{.Container.Name}} started' | mail -s 'Container Started' admin@example.com"
      chadburn.job-lifecycle.notify-on-start.event-type: "start"
      chadburn.job-lifecycle.cleanup-on-stop.command: "echo 'Container {{.Container.Name}} stopped' | mail -s 'Container Stopped' admin@example.com"
      chadburn.job-lifecycle.cleanup-on-stop.event-type: "stop"
```

#### Parameters

- `container`: The name or ID of the container to monitor.
- `event-type`: The type of event to trigger on. Valid values are `start` and `stop`.
- `command`: The command to execute when the event occurs.

#### Variables

You can use the following variables in your commands:

- `{{.Container.Name}}`: The name of the container.
- `{{.Container.ID}}`: The ID of the container.

## Job-exec

This job is executed inside a running container. Similar to `docker exec`

### Parameters

- **Schedule** *
  - *description*: When the job should be executed. E.g. every 10 seconds or every night at 1 AM.
  - *value*: String, see [Scheduling format](https://godoc.org/github.com/robfig/cron) of the Go implementation of `cron`. E.g. `@every 10s` or `0 0 1 * * *` (every night at 1 AM). **Note**: the format starts with seconds, instead of minutes.
  - *default*: Required field, no default.
- **Command** *
  - *description*: Command you want to run inside the container.
  - *value*: String, e.g. `touch /tmp/example`
  - *default*: Required field, no default.
- **Container** *
  - *description*: Name of the container you want to execute the command in.
  - *value*: String, e.g. `nginx-proxy`
  - *default*: Required field, no default.
- **User**
  - *description*: User as which the command should be executed, similar to `docker exec --user <user>`
  - *value*: String, e.g. `www-data`
  - *default*: `root`
- **tty**
  - *description*: Allocate a pseudo-tty, similar to `docker exec -t`. See this [Stack Overflow answer](https://stackoverflow.com/questions/30137135/confused-about-docker-t-option-to-allocate-a-pseudo-tty) for more info.
  - *value*: Boolean, either `false` or `true`
  - *default*: `false`
- **Workdir**
  - *description*: Working directory in which the command is executed, similar to `docker exec --workdir <dir>`
  - *value*: String, e.g. `/app`
  - *default*: Container's default working directory
  
### INI-file example

```ini
[job-exec "flush-nginx-logs"]
schedule = @hourly
container = nginx-proxy
command = /bin/bash /flush-logs.sh
user = www-data
tty = false
workdir = /var/log/nginx
```

### Docker labels example

```sh
docker run -it --rm \
    --label chadburn.enabled=true \
    --label chadburn.job-exec.flush-nginx-logs.schedule="@hourly" \
    --label chadburn.job-exec.flush-nginx-logs.command="/bin/bash /flush-logs.sh" \
    --label chadburn.job-exec.flush-nginx-logs.user="www-data" \
    --label chadburn.job-exec.flush-nginx-logs.tty="false" \
    --label chadburn.job-exec.flush-nginx-logs.workdir="/var/log/nginx" \
        nginx
```

## Job-run

This job can be used in 2 situations:

1. To run a command inside of a new container, using a specific image. Similar to `docker run`
1. To start a stopped container, similar to `docker start`

### Parameters

- **Schedule** * (1,2)
  - *description*: When the job should be executed. E.g. every 10 seconds or every night at 1 AM.
  - *value*: String, see [Scheduling format](https://godoc.org/github.com/robfig/cron) of the Go implementation of `cron`. E.g. `@every 10s` or `0 0 1 * * *` (every night at 1 AM). **Note**: the format starts with seconds, instead of minutes.
  - *default*: Required field, no default.
- **Command** (1)
  - *description*: Command you want to run inside the container.
  - *value*: String, e.g. `touch /tmp/example`
  - *default*: Default container command
- **Image** (1)
  - *description*: Image you want to use for the job.
  - *value*: String, e.g. `nginx:latest`
  - *default*: No default. If left blank, Chadburn assumes you will specify a container to start (situation 2).
- **User** (1)
  - *description*: User as which the command should be executed, similar to `docker run --user <user>`
  - *value*: String, e.g. `www-data`
  - *default*: `root`
- **Network** (1)
  - *description*: Connect the container to this network
  - *value*: String, e.g. `backend-proxy`
  - *default*: Optional field, no default.
- **Delete** (1)
  - *description*: Delete the container after the job is finished. Similar to `docker run --rm`
  - *value*: Boolean, either `true` or `false`
  - *default*: `true`
- **Container** (2)
  - *description*: Name of the container you want to start.
  - *value*: String, e.g. `nginx-proxy`
  - *default*: Required field in case parameter `image` is not specified, no default.
- **tty** (1,2)
  - *description*: Allocate a pseudo-tty, similar to `docker exec -t`. See this [Stack Overflow answer](https://stackoverflow.com/questions/30137135/confused-about-docker-t-option-to-allocate-a-pseudo-tty) for more info.
  - *value*: Boolean, either `true` or `false`
  - *default*: `false`
- **Volume**
  - *description*: Mount host machine directory into container as a [bind mount](https://docs.docker.com/storage/bind-mounts/#start-a-container-with-a-bind-mount)
  - *value*: Same format as used with `-v` flag within `docker run`. For example: `/tmp/test:/tmp/test:ro`
    - **INI config**: `Volume` setting can be provided multiple times for multiple mounts.
    - **Labels config**: multiple mounts has to be provided as JSON array: `["/test/tmp:/test/tmp:ro", "/test/tmp:/test/tmp:rw"]`
  - *default*: Optional field, no default.
  
### INI-file example

```ini
[job-run "print-write-date"]
schedule = @every 5s
image = alpine:latest
command = sh -c 'date | tee -a /tmp/test/date'
volume = /tmp/test:/tmp/test:rw
```

Then you can check output in host machine file `/tmp/test/date`

### Docker Compose example

When using `job-run` with Docker Compose, it's important to understand that this job type is designed to:
1. Run a command in a **new** container using a specific image, or
2. Start an **existing stopped** container

It is **not** designed to run commands inside already running containers (use `job-exec` for that).

#### Example 1: Running a command in a new container

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    command: daemon

  # This container doesn't need to be running for the job to work
  # It just defines the job configuration
  job-config:
    image: alpine:latest
    labels:
      chadburn.enabled: "true"
      chadburn.job-run.periodic-task.schedule: "@every 5m"
      chadburn.job-run.periodic-task.image: "alpine:latest"
      chadburn.job-run.periodic-task.command: "sh -c 'date | tee -a /tmp/log'"
      chadburn.job-run.periodic-task.volume: "/tmp:/tmp:rw"
```

In this example, Chadburn will create a new Alpine container every 5 minutes to run the date command.

#### Example 2: Starting an existing container periodically

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    command: daemon

  periodic-container:
    image: alpine:latest
    command: "sh -c 'date | tee -a /tmp/log'"
    labels:
      chadburn.job-run.restart-me.schedule: "@every 30m"
      chadburn.job-run.restart-me.container: "periodic-container"
```

In this example, the `periodic-container` will be started every 30 minutes. The container will execute its defined command each time it starts.

#### Common Mistake

A common mistake is trying to use `job-run` to execute commands inside a running container:

```yaml
# This will NOT work as expected
services:
  alpine:
    image: alpine
    # This container will start and stay running
    command: "tail -f /dev/null"
    labels:
      chadburn.enabled: "true"
      chadburn.job-run.datecron.schedule: "@every 5s"
      chadburn.job-run.datecron.command: "uname -a"
```

This configuration will not run `uname -a` every 5 seconds inside the running Alpine container. Instead, use `job-exec` for that purpose:

```yaml
# Correct way to run commands in an existing container
services:
  alpine:
    image: alpine
    command: "tail -f /dev/null"
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.datecron.schedule: "@every 5s"
      chadburn.job-exec.datecron.command: "uname -a"
```

### Running Chadburn on Docker example

```sh
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
        premoewb/chadburn:latest daemon 
```

## Job-local

Runs the command on the host running Chadburn.

**Note**: In case Chadburn is running inside a container, the command is executed inside the container. Not on the Docker host.

### Parameters

- **Schedule** *
  - *description*: When the job should be executed. E.g. every 10 seconds or every night at 1 AM.
  - *value*: String, see [Scheduling format](https://godoc.org/github.com/robfig/cron) of the Go implementation of `cron`. E.g. `@every 10s` or `0 0 1 * * *` (every night at 1 AM). **Note**: the format starts with seconds, instead of minutes.
  - *default*: Required field, no default.
- **Command** *
  - *description*: Command you want to run on the host.
  - *value*: String, e.g. `touch test.txt`
  - *default*: Required field, no default.
- **Dir**
  - *description*: Base directory to execute the command.
  - *value*: String, e.g. `/tmp/sandbox/`
  - *default*: Current directory
- **Environment** (Broken?)
  - *description*: List of environment variables
  - *value*: String, e.g. `FILE=test.txt`
  - *default*: Optional field, no default.

### INI-file example

```ini
[job-local "create-file"]
schedule = @every 15s
command = touch test.txt
dir = /tmp/
```

## Job-service-run

This job can be used to:

- To run a command inside a new "run-once" service, for running inside a swarm.

### Parameters

- **Schedule** * (1,2)
  - *description*: When the job should be executed. E.g. every 10 seconds or every night at 1 AM.
  - *value*: String, see [Scheduling format](https://godoc.org/github.com/robfig/cron) of the Go implementation of `cron`. E.g. `@every 10s` or `0 0 1 * * *` (every night at 1 AM). **Note**: the format starts with seconds, instead of minutes.
  - *default*: Required field, no default.
- **Command** (1, 2)
  - *description*: Command you want to run inside the container.
  - *value*: String, e.g. `touch /tmp/example`
  - *default*: Default container command
- **Image** * (1)
  - *description*: Image you want to use for the job.
  - *value*: String, e.g. `nginx:latest`
  - *default*: No default. If left blank, Chadburn assumes you will specify a container to start (situation 2).
- **Network** (1)
  - *description*: Connect the container to this network
  - *value*: String, e.g. `backend-proxy`
  - *default*: Optional field, no default.
- **delete** (1)
  - *description*: Delete the container after the job is finished.
  - *value*: Boolean, either `true` or `false`
  - *default*: `true`
- **User** (1,2)
  - *description*: User as which the command should be executed.
  - *value*: String, e.g. `www-data`
  - *default*: `root`
- **tty** (1,2)
  - *description*: Allocate a pseudo-tty, similar to `docker exec -t`. See this [Stack Overflow answer](https://stackoverflow.com/questions/30137135/confused-about-docker-t-option-to-allocate-a-pseudo-tty) for more info.
  - *value*: Boolean, either `true` or `false`
  - *default*: `false`
  
### INI-file example

```ini
[job-service-run "service-executed-on-new-container"]
schedule = 0,20,40 * * * *
image = ubuntu
network = swarm_network
command =  touch /tmp/example
```

## Variable Substitution

Chadburn supports variable substitution in job commands, allowing you to reference container information dynamically. This is particularly useful when you need to interact with containers without hardcoding their names or IDs.

### Available Variables

- **{{.Container.Name}}** - The name of the container
- **{{.Container.ID}}** - The ID of the container

### Examples

```ini
[job-exec "restart-container"]
schedule = @daily
container = my-container
command = docker restart {{.Container.Name}}

[job-local "log-container-id"]
schedule = @hourly
command = echo "Container ID is {{.Container.ID}}" >> /var/log/containers.log
```

With Docker labels:

```sh
docker run -it --rm \
    --label chadburn.enabled=true \
    --label chadburn.job-exec.restart-self.schedule="@daily" \
    --label chadburn.job-exec.restart-self.command="docker restart {{.Container.Name}}" \
        nginx
```

This feature allows for more flexible and reusable job configurations, especially when working with dynamic container environments where container names or IDs might change.
