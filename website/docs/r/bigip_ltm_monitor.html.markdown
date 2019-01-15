---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_monitor"
sidebar_current: "docs-bigip-resource-monitor-x"
description: |-
    Provides details about bigip_ltm_monitor resource
---

# bigip\_ltm\_monitor

`bigip_ltm_monitor` Configures a custom monitor for use by health checks.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage


```hcl
resource "bigip_ltm_monitor" "monitor" {
  name = "/Common/terraform_monitor"
  parent = "/Common/http"
  send = "GET /some/path\r\n"
  timeout = "999"
  interval = "999"
  destination = "1.2.3.4:1234"
}

```      

## Argument Reference

* `name` (Required) Name of the monitor

* `parent` - (Required) Existing LTM monitor to inherit from

* `interval` - (Optional) Check interval in seconds

* `timeout` - (Optional) Timeout in seconds

* `send` - (Optional) Request string to send

* `receive` - (Optional) Expected response string

* `receive_disable` - (Optional)

* `reverse`  - (Optional)

* `transparent` - (Optional)

* `manual_resume` - (Optional)

* `ip_dscp` - (Optional)

* `time_until_up` - (Optional)

* `destination` - (Optional) Specify an alias address for monitoring

* `compatibility` -  (Optional) Specifies, when enabled, that the SSL options setting (in OpenSSL) is set to ALL. Accepts 'enabled' or 'disabled' values, the default value is 'enabled'.
