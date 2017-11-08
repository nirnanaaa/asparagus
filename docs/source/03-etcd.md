# ETCD Source Provider

Asparagus has built-in support for high-availability and distributed execution
configuration.

ETCD is used as its backend, to ensure all cluster members share the same information.
It can be used alonside other source providers like the crontab provider, for example.


## Installation

This module is included in the core of asparagus. As such, you can just go to your
`asparagus.conf` and enable it.

```toml
[scheduler.provider-etcd]
  # Required: Enable the ETCD Provider. Without this flag all following settings would be
  # ignored.
  enabled = true

  # Optional: Define a registry url.
  # Alternatively you can enable dns discovery and skip this entry.
  registry-url = ["http://127.0.0.1:4001"]

  # Optional: Enable DNS Discovery itself
  # Guide:
  # https://coreos.com/etcd/docs/latest/v2/clustering.html#dns-discovery
  discovery-enabled = false

  # Optional: Define a discovery domain (leave the SRV record parts out)
  discovery-domain = "example.com"

  # Optional: Define a CA-Cert, used to connect to the ETCD nodes.
  connect-ca-cert = ""

  # Optional: define an alternate folder location, where your cronjobs will be
  # stored.
  source-folder = "/cron/Jobs"
```

## Job Configuration

All jobs must be stored under an unique key inside the `source-folder`, configured
inside your `asparagus.conf`.


### Configuration

```json
{
  "executor": "string",
  "schedule":"cron-expr",
  "execution_config": {},
  "is_running": false,
  "last_run_at": "rfc3339 date format",
  "name": "optional string"
}
```

### Example

```bash
etcdctl set /cron/Jobs/TestJob '{"executor": "http", "schedule":"* * * * *", "execution_config":{"URL": "https://httpbin.org/get", "Method": "GET"}}'
```
