# Asparagus

[![CircleCI](https://circleci.com/gh/nirnanaaa/asparagus.svg?style=svg)](https://circleci.com/gh/nirnanaaa/asparagus)
[![Maintainability](https://api.codeclimate.com/v1/badges/ef3bb9a0dbc5b1994a9d/maintainability)](https://codeclimate.com/github/nirnanaaa/asparagus/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/ef3bb9a0dbc5b1994a9d/test_coverage)](https://codeclimate.com/github/nirnanaaa/asparagus/test_coverage)

A simple and pluggable cron scheduler for distributed systems.

# Usage

```bash

# Crontab Backend
docker pull nirnanaaa/asparagus

# Optional: Spin up an etcd instance
docker run \
  --rm -ti \
  -v $(pwd)/example:/app \
  -e ASPARAGUS_SCHEDULER_PROVIDER_CRONTAB_SOURCE=/app/crontab \
  nirnanaaa/asparagus
```

# How it works

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

# Configuration

## Global configuration

Configuration is done using environment variables or a central config file. You can generate a sample configuration file by using the following command:

```
docker run --rm -ti nirnanaaa/asparagus config
```

Now you can start your asparagus instance using the provided volume: `/etc/asparagus` using the `asparagus.conf` filename:

```
docker run --rm -ti -v $(pwd)/conf:/etc/asparagus nirnanaaa/asparagus config
```

Alternatively you can use prefixed environment variables for each configuration variable, like so:

```
ASPARAGUS_<SECTION>_<KEY>=<VALUE>
```
# Crontab/Source providers
## local crontab

See the `example/crontab` file for a reference. The sample config inside the docs folder should guide you through setting everyting up.


## SQL (postgres/mysql/sqlite..., tbd)

## ETCD - TBD Reimplementation

# Execution provider

## HTTP
## Kafka (tbd)
## SQS (tbd)


# License

Copyright 2017 Florian Kasper

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
