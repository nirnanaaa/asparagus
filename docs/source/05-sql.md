# SQL Source Provider

Asparagus may be used alongside a database service. Supported databases are:

- sqlite3
- mysql

Asparagus will create a database table for you, providing you specify a correct database.

## Installation

This module is included in the core of asparagus. As such, you can just go to your
`asparagus.conf` and enable it.

```toml
[scheduler.provider-sql]
  # Required: Enable the SQL provider.
  enabled = true

  # for sqlite3
  driver = "sqlite3"
  uri = "/tmp/test.db"

  # for mysql
  driver = "mysql"

  # make sure to suffix ?parseTime=true otherwise asparagus will complain about string to time conversion.
  uri = "root:pass@tcp(127.0.0.1:3306)/database?parseTime=true"
```

## Job Configuration

All jobs must be stored under an unique name inside the database.


### Example

```bash
insert into cronjobs values ('test', 'test', '0', '* * * * *', NOW(), '', 'http', '{\"URL\": \"https://httpbin.org/get\", \"Method\": \"GET\"}', '', '10', '0');
```
