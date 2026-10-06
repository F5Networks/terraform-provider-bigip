---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_syslog"
subcategory: "System"
description: |-
  Provides details about bigip_sys_syslog resource
---

# bigip\_sys\_syslog

`bigip_sys_syslog` manages the BIG-IP system's syslog settings: per-facility log levels and the list of remote syslog servers the system forwards messages to.

This is a singleton resource -- BIG-IP has exactly one syslog configuration, so `Delete` resets it to BIG-IP's factory defaults rather than removing an object.

## Example Usage

```hcl
resource "bigip_sys_syslog" "example" {
  auth_priv_from = "notice"
  auth_priv_to   = "emerg"
  console_log    = "enabled"

  remote_servers {
    name        = "splunk-1"
    host        = "10.10.10.10"
    remote_port = 514
  }

  remote_servers {
    name        = "splunk-2"
    host        = "10.10.10.11"
    remote_port = 6514
  }
}
```

## Argument Reference

* `auth_priv_from` - (Optional, type `string`, Computed) Lowest level of messages from the auth/priv facility to include in the log.

* `auth_priv_to` - (Optional, type `string`, Computed) Highest level of messages from the auth/priv facility to include in the log.

* `console_log` - (Optional, type `string`, Computed) Whether logs are provided to the console. `enabled` or `disabled`.

* `cron_from` / `cron_to` - (Optional, type `string`, Computed) Log level range for the cron facility.

* `daemon_from` / `daemon_to` - (Optional, type `string`, Computed) Log level range for the daemon facility.

* `iso_date` - (Optional, type `string`, Computed) Whether the date included in a log message conforms to the ISO 8601 standard. `enabled` or `disabled`.

* `kern_from` / `kern_to` - (Optional, type `string`, Computed) Log level range for the kernel facility.

* `local6_from` / `local6_to` - (Optional, type `string`, Computed) Log level range for the local6 facility.

* `mail_from` / `mail_to` - (Optional, type `string`, Computed) Log level range for the mail facility.

* `messages_from` / `messages_to` - (Optional, type `string`, Computed) Log level range for the user/messages facility.

* `user_log_from` / `user_log_to` - (Optional, type `string`, Computed) Log level range for the user-log facility.

* `remote_servers` - (Optional, type `list`) One or more remote syslog server blocks. Each block supports:
  * `name` - (Required, type `string`) Name of the remote syslog server entry.
  * `host` - (Required, type `string`) IP address or hostname of the remote syslog server.
  * `local_ip` - (Optional, type `string`) Local IP address the BIG-IP system uses when sending to this server.
  * `remote_port` - (Optional, type `int`, Default `514`) Port on the remote syslog server.

## Import

An existing syslog configuration can be imported with:

```sh
terraform import bigip_sys_syslog.example syslog
```
