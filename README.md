# Chadburn - a job scheduler written in Go [![GitHub version](https://badge.fury.io/gh/PremoWeb%2FChadburn.svg)](https://github.com/PremoWeb/Chadburn/releases) ![Test](https://github.com/PremoWeb/Chadburn/workflows/Test/badge.svg)

**Chadburn** is a modern and low footprint job scheduler for __docker__ environments, built on Go. Chadburn aims to be a replacement for the old fashioned [cron](https://en.wikipedia.org/wiki/Cron).

*** SPECIAL NOTE ***

Chadburn is a new project based on the previous and continuous work incorporated into Ofelia and a fork of Ofelia provided by @rdelcorro, of which Chadburn was forked from. This project was started as a result of needing a version of Ofelia that incorporated the following fixes:

- Update tasks if docker containers are started, stopped, restarted, or changed
- Do not require a dummy task on the Ofelia container just to use Chadburn.
- Support INI and docker labels at the same time. The configs will simply be merged.
- Do not require a restart in order to pick up new or remove tasks.

PremoWeb will be responsive to addressing issues raised in this project and will also be monitoring issues in the original Ofelia source code repository and applying changes that should be reflected in Chadburn.
### Why Chadburn?

It has been a long time since [`cron`](https://en.wikipedia.org/wiki/Cron) was released, actually more than 28 years. The world has changed a lot and especially since the `Docker` revolution. **Vixie's cron** works great but it's not extensible and it's hard to debug when something goes wrong.

Many solutions are available: ready to go containerized `crons`, wrappers for your commands, etc. but in the end simple tasks become complex.

### How?

The main feature of **Chadburn** is the ability to execute commands directly on Docker containers. Using Docker's API Chadburn emulates the behavior of [`exec`](https://docs.docker.com/reference/commandline/exec/), being able to run a command inside of a running container. Also you can run the command in a new container destroying it at the end of the execution.

## Configuration

A wiki is being written to document how to use Chadburn. Caprover users can use a One Click App (coming soon) to deploy and implement scheduled jobs using Service Label Overrides.

For everyone else, here's the general approach to use Chadburn:

### Jobs

[Scheduling format](https://godoc.org/github.com/robfig/cron) is the same as the Go implementation of `cron`. E.g. `@every 10s` or `0 0 1 * * *` (every night at 1 AM).

**Note**: the format starts with seconds, instead of minutes.

you can configure four different kind of jobs:

- `job-exec`: this job is executed inside of a running container.
- `job-run`: runs a command inside of a new container, using a specific image.
- `job-local`: runs the command inside of the host running Chadburn.
- `job-service-run`: runs the command inside a new "run-once" service, for running inside a swarm

See [Jobs reference documentation](docs/jobs.md) for all available parameters.

#### INI-style config

Run with `chadburn daemon --config=/path/to/config.ini`

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
command =  touch /tmp/example
```

#### Docker labels configurations

In order to use this type of configurations, Chadburn need access to docker socket.

```sh
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
        premoweb/chadburn:latest daemon --docker
```

Labels format: `chadburn.<JOB_TYPE>.<JOB_NAME>.<JOB_PARAMETER>=<PARAMETER_VALUE>.
This type of configuration supports all the capabilities provided by INI files.

Also, it is possible to configure `job-exec` by setting labels configurations on the target container. To do that, additional label `chadburn.enabled=true` need to be present on the target container.

For example, we want `chadburn` to execute `uname -a` command in the existing container called `my_nginx`.
To do that, we need to we need to start `my_nginx` container with next configurations:

```sh
docker run -it --rm \
    --label chadburn.enabled=true \
    --label chadburn.job-exec.test-exec-job.schedule="@every 5s" \
    --label chadburn.job-exec.test-exec-job.command="uname -a" \
        nginx
```

Now if we start `chadburn` container with the command provided above, it will execute the task:

- Exec  - `uname -a`

Or with docker-compose:

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    depends_on:
      - nginx
    command: daemon --docker
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro

  nginx:
    image: nginx
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.datecron.schedule: "@every 5s"
      chadburn.job-exec.datecron.command: "uname -a"
```

#### Dynamic docker configuration

You can start Chadburn in its own container or on the host itself, and it will magically pick up any container that starts, stops or is modified on the fly.
In order to achieve this, you simply have to use docker containers with the labels described above and let Chadburn take care of the rest. 

#### Hybrid configuration (INI files + Docker)

You can specify part of the configuration on the INI files, such as globals for the middlewares or even declare tasks in there but also merge them with docker.
The docker labels will be parsed, added and removed on the fly but also, the file config can be used to execute tasks that are not possible using just docker labels 
such as:

- job-local
- job-run

**Use the INI file to:**

- Configure the slack or other middleware integration
- Configure any global setting
- Create a job-run so it executes on a new container each time

```ini
[global]
slack-webhook = https://myhook.com/auth

[job-run "job-executed-on-new-container"]
schedule = @hourly
image = ubuntu:latest
command = touch /tmp/example
```

**Use docker to:**

```sh
docker run -it --rm \
    --label chadburn.enabled=true \
    --label chadburn.job-exec.test-exec-job.schedule="@every 5s" \
    --label chadburn.job-exec.test-exec-job.command="uname -a" \
        nginx
```

### Logging
**Chadburn** comes with three different logging drivers that can be configured in the `[global]` section:
- `mail` to send mails
- `save` to save structured execution reports to a directory
- `slack` to send messages via a slack webhook

#### Options
- `smtp-host` - address of the SMTP server.
- `smtp-port` - port number of the SMTP server.
- `smtp-user` - user name used to connect to the SMTP server.
- `smtp-password` - password used to connect to the SMTP server.
- `email-to` - mail address of the receiver of the mail.
- `email-from` - mail address of the sender of the mail.
- `mail-only-on-error` - only send a mail if the execution was not successful.

- `save-folder` - directory in which the reports shall be written.
- `save-only-on-error` - only save a report if the execution was not successful.

- `slack-webhook` - URL of the slack webhook.
- `slack-only-on-error` - only send a slack message if the execution was not successful.

### Overlap
**Chadburn** can prevent that a job is run twice in parallel (e.g. if the first execution didn't complete before a second execution was scheduled. If a job has the option `no-overlap` set, it will not be run concurrently. 

## Installation

The easiest way to deploy **Chadburn** is using *Docker*. See examples above.

If don't want to run **Chadburn** using our *Docker* image you can download a binary from [releases](https://github.com/PremoWeb/Chadburn/releases) page.

### Thank You to team Ofelia and it's contributors.

A special thanks to [@rdelcorro](https://github.com/rdelcorro) for the work in fixing the issues referenced in this pull request https://github.com/mcuadros/ofelia/pull/137, despite this pull request having been ignored for 30 days. PremoWeb aims to ensure that open software is continously improve and will remain responsive to raised issues and pull requests.

Much thanks to the original work that went into Ofelia by it's author and contributors.