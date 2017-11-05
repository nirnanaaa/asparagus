```toml
[meta]
  # Enable debug mode. This enables some additional features, like PPROF web
  # interface.
  debug = false

  # Enable logging AT ALL
  logging-enabled = true

  # Set the log level of the application
  log-level = "warn"

[metrics]
  # Enable metrics exporting in general
  enable-metrics = false

  # Slack metrics exporter
  [metrics.slack]
    enabled = false
    webhook-url = ""
    message-format = ""
    reporting-interval = "24h0m0s"

  # Console metrics exporter
  [metrics.console]
    enabled = false
    reporting-interval = "24h0m0s"

  # Cloudwatch metrics exporter
  [metrics.cloudwatch]
    enabled = false
    message-format = ""
    reporting-interval = "24h0m0s"
    namespace = ""
    reset-counters-after-report = false

    [metrics.cloudwatch.static-dimensions]
    # variables that will be transferred to cloudwatch.
    # you can use any variable you like.
    # just make sure to use the
    # key = "stringvalue" format
    environment = "stage"

# Task scheduler and dispatcher configuration
[scheduler]
  # Enable or disable the task scheduler. Could be useful for setting up the system
  enabled = true

  # How often should the scheduler check for new jobs to be executed
  tick-duration = "1s"

  # Log when new tasks were inserted/updated
  log-task-detection = false

  # How many workers should run to execute work - one worker can execute one job at a time
  # so make sure to scale them properly.
  num-workers = 10

  # The HTTP Executor (cronjob target)
  [scheduler.executor-http]
    # Enable it
    enabled = true

    # Debug the HTTP request/response cycle. Could be useful for detecting job errors.
    # Will log an output like `curl -v` into stdout.
    debug-response = false

    # Log a success/error message with the HTTP status code into console.
    # this is unaffected by `debug-response`
    log-http-status = false

    # Sign the requests using JWT tokens inside the `Authorization` header.
    # use with caution. This will probably be exported into a seperate executor
    # or middleware.
    use-jwt-signing = true
    jwt-issuer = "issuer"
    jwt-expires = "10m0s"
    jwt-subject = "0"
    jwt-secret = ""

  # The Crontab source provider (crontab source)
  [scheduler.provider-crontab]
    enabled = true
    source = "/etc/crontab"
```
