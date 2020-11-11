---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_monitor"
sidebar_current: "docs-bigip-datasource-monitor-x"
description: |-
    Provides details about bigip_ltm_monitor data source
---

# bigip\_ltm\_monitor

Use this data source (`bigip_ltm_monitor`) to get the ltm monitor details available on BIG-IP
 
 
## Example Usage
```hcl

data "bigip_ltm_monitor" "Monitor-TC1" {
  name = "test-monitor"
  partition = "Common"
}

```      

## Argument Reference

* `name` - (Required) Name of the ltm monitor

* `partition` - (Required) partition of the ltm monitor


## Attributes Reference

Additionally, the following attributes are exported:

* `destination` - id will be full path name of ltm monitor.

* `interval` - id will be full path name of ltm monitor.

* `ip_dscp` - id will be full path name of ltm monitor.

* `manual_resume` - id will be full path name of ltm monitor.

* `reverse` - id will be full path name of ltm monitor.

* `transparent` - id will be full path name of ltm monitor.

* `adaptive_limit` - id will be full path name of ltm monitor.

* `adaptive` - id will be full path name of ltm monitor.
