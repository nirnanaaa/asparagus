# Asparagus

[![CircleCI](https://circleci.com/gh/nirnanaaa/asparagus.svg?style=svg)](https://circleci.com/gh/nirnanaaa/asparagus)
[![Maintainability](https://api.codeclimate.com/v1/badges/ef3bb9a0dbc5b1994a9d/maintainability)](https://codeclimate.com/github/nirnanaaa/asparagus/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/ef3bb9a0dbc5b1994a9d/test_coverage)](https://codeclimate.com/github/nirnanaaa/asparagus/test_coverage)

Simple cron scheudler for distributed systems.

# Usage

```bash
docker pull nirnanaaa/asparagus

# Optional: Spin up an etcd instance
# docker run --name etcd -d elcolio/etcd:latest

docker run --rm -ti -e ASPARAGUS_ETCD_REGISTRY_URL=http://etcd:4001 --link etcd:etcd nirnanaaa/asparagus
```

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

## Cronjob configuration

Each cronjob follows a pre-defined json schema. It uses a regular cron expression for execution:


```json
{
  "Name":"string",
  "Expression":"0 02,03,23 * * *",
  "LastRunAt":"rfc3339 formatted datetime",
  "AfterTask":"string - following the \"Name\" variable",
  "URI":"string",
  "URIIsAbsolute":true,
}
```

To see all possible configuration options take a look inside the `scheduler/definition.go` file.

Those cronjob should be placed inside the configured cronjob folder inside etcd (default: `/cron/Jobs/<name>`)

## Scheduler configuration
The scheduler itself has hot-configuration mode, that can be used to alter tick interval and a jwt secret to sign your requests with (ROADMAP: configureable). Its location can be specified inside the `/cron/Config` (ROADMAP) key inside etcd.

Schema:

```json
{
  "SchedulerTickDuration":1000,
  "Enabled":true,
  "SecretKey":"string"
}
```

# How it works


# License

Copyright 2017 Florian Kasper

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
