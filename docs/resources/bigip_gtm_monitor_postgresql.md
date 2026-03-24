# bigip_gtm_monitor_postgresql Resource

Provides a BIG-IP GTM (Global Traffic Manager) PostgreSQL Monitor resource. This resource allows you to configure and manage GTM PostgreSQL health monitors on a BIG-IP system.

## Description

A GTM PostgreSQL monitor verifies PostgreSQL database services by connecting to a database and optionally executing a query and evaluating the response. PostgreSQL monitors support database-specific configuration including database name, username, password, and instance count.

## Example Usage

### Basic PostgreSQL Monitor

```hcl
resource "bigip_gtm_monitor_postgresql" "example" {
  name = "/Common/my_postgresql_monitor"
}
```

### PostgreSQL Monitor with Authentication

```hcl
resource "bigip_gtm_monitor_postgresql" "advanced" {
  name                 = "/Common/my_postgresql_monitor"
  defaults_from        = "/Common/postgresql"
  destination          = "*:5432"
  interval             = 10
  timeout              = 60
  probe_timeout        = 3
  ignore_down_response = "disabled"
  database             = "mydb"
  username             = "monitor_user"
  password             = "monitor_pass"
  receive              = "SELECT"
  debug                = "no"
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required, String) The full path name of the GTM PostgreSQL monitor (e.g., `/Common/my_postgresql_monitor`). Forces new resource.

### Optional Arguments

#### General Settings

* `defaults_from` - (Optional, String) Specifies the parent monitor from which this monitor inherits settings. Default: `/Common/postgresql`.
* `destination` - (Optional, String) Specifies the IP address and service port of the resource being monitored. Format: `ip:port`. Default: `*:*`.
* `interval` - (Optional, Integer) Specifies, in seconds, the frequency at which the system issues the monitor check. Default: `30`.
* `timeout` - (Optional, Integer) Specifies the number of seconds the target has in which to respond to the monitor request. Default: `91`.
* `probe_timeout` - (Optional, Integer) Specifies the number of seconds after which the system times out the probe request. Default: `5`.
* `ignore_down_response` - (Optional, String) Specifies whether the monitor ignores a down response from the system it is monitoring. Valid values: `enabled`, `disabled`. Default: `disabled`.

#### Database Settings

* `database` - (Optional, String) Specifies the name of the database that the monitor tries to access.
* `username` - (Optional, String) Specifies the user name if the monitored target requires authentication.
* `password` - (Optional, String, Sensitive) Specifies the password if the monitored target requires authentication.
* `receive` - (Optional, String) Specifies the text string that the monitor looks for in the returned resource.

#### Connection Settings

* `instance_count` - (Optional, String) Specifies the number of instances for which the system keeps a connection open. By default, when you assign instances of this monitor to a resource, the system keeps the connection to the database open. With this option you can assign multiple instances to the database while reducing the overhead that multiple open connections can cause. A value of `0` keeps the connection open for all instances. A value of `1` opens a new connection for each instance. Any other positive value keeps the connection open for that many instances.
* `debug` - (Optional, String) Specifies whether the monitor sends error messages and additional information to a log file created and labeled specifically for this monitor. Valid values: `yes`, `no`. Default: `no`.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The full path name of the GTM PostgreSQL monitor.

## Import

GTM PostgreSQL Monitor resources can be imported using the full path name:

```bash
terraform import bigip_gtm_monitor_postgresql.example /Common/my_postgresql_monitor
```
