---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_ntp"
subcategory: "System"
description: |-
  Provides details about bigip_sys_ntp data source
---

# bigip\_sys\_ntp

Use this data source (`bigip_sys_ntp`) to retrieve the BIG-IP system's NTP configuration (time servers and time zone). `sys/ntp` is a singleton object -- there is exactly one NTP configuration per BIG-IP system -- so this data source takes no arguments.

Use [`bigip_sys_ntp`](../resources/bigip_sys_ntp.md) (the resource) to manage this configuration. Note the resource's time-server argument is named `servers`; this data source intentionally exposes the same list as `ntp_servers` instead.

## Example Usage

```hcl
data "bigip_sys_ntp" "ntp" {}

output "ntp_servers" {
  value = data.bigip_sys_ntp.ntp.ntp_servers
}

output "timezone" {
  value = data.bigip_sys_ntp.ntp.timezone
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `description` - User-defined description of the system NTP configuration.

* `ntp_servers` - Time servers the system uses to update the system time.

* `timezone` - Time zone used for the system time.
