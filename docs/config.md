```toml
[meta]
  debug = true
  logging-enabled = true
  log-level = 5
  Version = ""

[metrics]
  enable-metrics = false
  [metrics.slack]
    enabled = false
    webook-url = ""
    message-format = ""
    reporting-interval = "24h0m0s"
  [metrics.console]
    enabled = false
    reporting-interval = "24h0m0s"
  [metrics.cloudwatch]
    enabled = false
    webook-url = ""
    message-format = ""
    reporting-interval = "24h0m0s"
    namespace = ""
    reset-counters-after-report = false
    [metrics.cloudwatch.static-dimensions]

[scheduler]
  enabled = true
  tick-duration = "1s"
  api-fallback-domain = "example.com"
  cronjob-base-folder = "/cron/Jobs"
  hot-config-path = "/cron/Config"
  log-task-detection = true
  num-workers = 10
  [scheduler.executor-http]
    enabled = true
    debug-response = true
    log-http-status = true
    use-jwt-signing = false
    jwt-issuer = "issuer"
    jwt-expires = "10m0s"
    jwt-subject = "0"
    jwt-secret = ""
  [scheduler.provider-crontab]
    enabled = true
    source = "/Users/fkasper/.asparagus/crontab"

[etcd]
  service-name = "DummyService"
  registry-url = ["http://127.0.0.1:4001"]
  binding-port = 9092
  local = false
  discovery-domain = "services.example.com"
  use-dns-discovery = false
  ca-cert = "/etc/ssl/certs/ca.crt"
  use-ssl = false
```