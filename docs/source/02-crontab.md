# Crontab Source Provider

Asparagus uses an alternated crontab syntax for its configuration. It consists
of three "chunks" of information:

```
[SCHEDULE] [TASK CONFIGURATION] [EXECUTOR CONFIGURATION]
```

- Schedule: Just like with regular cron you can define a schedule here. To read
  more, please read the following article: [Article](https://help.ubuntu.com/community/CronHowto)
  It also has support for `@hourly` `@daily` and so on, except for `@reboot`
- Task configuration: Here you can define some basic parameters, like dependencies,
  naming and if the cronjob is enabled. Use the following variables inside square
  backets like so: `[Name=MyName]` alternatively you can use space separated values
  as well: `[Name="My Name"]`. To assign multiple variables use a space between
  the definitions: `[Name="My Name" Executor="http" AfterTask="My Name1"]`. The
  following variables can be assigned:

  - `Name`: defines a name for the cronjob. *Required*. Is used for the AfterTask
    variable
  - `Executor`: defines which executor should be used. *Required*. See [this](../executors.md)
    for a list of valid executors
  - `AfterTask`: defines a dependency for this task. Use the `Name` of the depending task.
- Executor configuration: Just like with the task configuration you can define settings
  that will be passed down to the executor. In case of using the HTTP executor, sample settings
  would be a `URL` or `Method`: `[URL="https://httpbin.org/get" Method="GET"]`.


# Examples

To see a valid crontab checkout [THIS](../../example/crontab) file.


- Execute a cronjob every minute, using the HTTP executor:
  ```
  * * * * * [Name="HTTP Test" Executor=http] [URL="https://httpbin.org/get" Method=get]
  ```
- Execute a cronjob every day, using the HTTP executor:
  ```
  @daily [Name="HTTP Test" Executor=http] [URL="https://httpbin.org/get" Method=get]
  ```
