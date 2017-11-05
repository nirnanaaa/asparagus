# Asparagus

[![CircleCI](https://circleci.com/gh/nirnanaaa/asparagus.svg?style=svg)](https://circleci.com/gh/nirnanaaa/asparagus)
[![Maintainability](https://api.codeclimate.com/v1/badges/ef3bb9a0dbc5b1994a9d/maintainability)](https://codeclimate.com/github/nirnanaaa/asparagus/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/ef3bb9a0dbc5b1994a9d/test_coverage)](https://codeclimate.com/github/nirnanaaa/asparagus/test_coverage)

## A simple and pluggable cron scheduler for distributed systems.

Asparagus is an open source **task scheduler** with **no external dependencies**. It's useful for
executing scheduled tasks on any environment.

## Features

- Variety of [Source Providers](./docs/source/01-base.md), ranging from standard crontab to etcd based backends.
- Multiple [Execution Providers](./docs/execution/01-base.md), including HTTP and local execution.
- Simple to install and manage, and easy to get started.
- Monitoring and reporting is a key priority.

## Installation

We recommend using Asparagus with our pre-built Docker image. Start Asparagus using:

```bash
# Crontab Backend
docker pull fkconsultin/asparagus

# Optional: Spin up an etcd instance
docker run \
  --rm -ti \
  -v $(pwd)/example:/app \
  -e ASPARAGUS_SCHEDULER_PROVIDER_CRONTAB_SOURCE=/app/crontab \
  fkconsultin/asparagus
```

## How it works

```
  -------------------                                         ----------------------
  | Source Provider | \                                     / | Execution Provider |
  -------------------  \                                   /  ----------------------
                        > - - - Asparagus Scheduler - - - <
  -------------------  /                |                  \  ----------------------
  | Source Provider | /                 |                   \ | Execution Provider |
  -------------------                   |                     ----------------------
                                ---------------------
                                | Metrics Reporting |
                                ---------------------
```

`Source Providers` provide configuration and crontab settings to asparagus.
Asparagus uses the configured `Execution Provider` to execute the actual cronjob
on the target system.

## Configuration

### Global configuration

Configuration is done using environment variables or a central config file. You can generate a sample configuration file by using the following command:

```
docker run --rm -ti fkconsultin/asparagus config > conf/asparagus.conf
```

Now you can start your asparagus instance using the provided volume: `/etc/asparagus` using the `asparagus.conf` filename:

```
docker run --rm -ti -v $(pwd)/conf:/etc/asparagus fkconsultin/asparagus
```

Alternatively you can use prefixed environment variables for each configuration variable, like so:

```
ASPARAGUS_<SECTION>_<KEY>=<VALUE>
```

Read the annotated documentation if you want to find out more: [here](./docs/config.md)

## Source Providers

Source providers define where cronjobs can be discovered. To find out more read [this](./docs/source/01-base.md) article.

## Execution Providers

Execution providers execute whatever was defined inside a cronjob. To find out more read [this](./docs/execution/01-base.md) article.

## Contributing

If you're feeling adventurous and want to contribute to Aspragus, see our [contributing doc](./CONTRIBUTING.md) for info on how to make feature requests, build from source, and run tests.
