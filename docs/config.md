```toml
[meta]
  # Enable debug mode
  debug = true

  # Enable logging in general
  logging-enabled = true

  # Log level configuration
  # 0 - Panic
  # 1 - Fatal
  # 2 - Error
  # 3 - Warn
  # 4 - Info
  # 5 - Debug
  log-level = 3

[metrics]
  # Enable metrics exporting in general
  enable-metrics = false

  # Slack reporter
  [metrics.slack]
    enabled = false

    # Slack webhook url
    webook-url = ""

    # Metrics message format. Use template/text syntax.
    #
    # Accessible variables:
    # .Metrics
    #   .Name
    #   .Value (counter: count, more tbd)
    message-format = '''
    Metrics export:
    {{range $index, $element := .Metrics}}
    {{range $element}}
      {{.Name}}: {{.Value}}
    {{end}}
    {{end}}
    '''
    # How often should messages be sent to slack?
    reporting-interval = 1d

  # Cloudwatch metrics reporter
  #
  # Make sure to have the correct IAM Instance policy/Credentials set up.
  #
  #
  [metrics.cloudwatch]
    enabled = false
    reporting-interval = 1m

[scheduler]
  tick-duration-ms = 1000
  api-fallback-domain = "example.com"
  cronjob-base-folder = "/cron/Jobs"

  # Path for changing configuration options on the fly inside the etcd adapter
  hot-config-path = "/cron/Config"

  # Log tasks upon startup
  log-task-detection = false
  [scheduler.http]
    # Dump HTTP Response - useful for debugging
    log-response-body = false

    # Log HTTP Status codes
    log-response-status = true

[etcd]
  # Registers itself inside the etcd service registry
  enable-registration = true

  # Uses service-name as etcd key prefix
  service-name = "DummyService"

  # Contact points for etcd
  registry-url = ["http://127.0.0.1:4001"]

  # Exposed, public port
  binding-port = 9092
  local = false

  # Domain, used for DNS based discovery.
  #
  # do not add the complete srv target path (_etcd-client-ssl._tcp.example.com => example.com)
  discovery-domain = "example.com"

  # Enable etcd server dns discovery. (SRV request)
  #
  # https://coreos.com/etcd/docs/latest/v2/clustering.html
  use-dns-discovery = false

  # SSL CA cert for self-signed ETCD certificates
  ca-cert = "/etc/ssl/certs/ca.crt"

  # Enable SSL mode (transparent and can be kept on even if not using it as long as the file under ca-cert exists)
  use-ssl = false
```
