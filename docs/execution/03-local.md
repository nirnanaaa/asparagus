# Local Execution Provider

The local provider is used to execute local commands on the machine, that runs asparagus.

## Configuration

By default the executor is disabled for security reasons. To enable it simply go
to your config and enable the `enabled` flag.

```
[scheduler.executor-local]
  enabled = false
  enable-output = false
```

# Examples

To see a valid crontab checkout [THIS](../../example/crontab) file.


- Execute a cronjob every minute, touching a file
  ```
  * * * * * [Name="Local Test" Executor=local] [touch /tmp/test.file]
  ```
