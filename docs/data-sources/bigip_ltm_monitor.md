---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_monitor"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_monitor data source
---

# bigip\_ltm\_monitor

Use this data source (`bigip_ltm_monitor`) to get the ltm monitor details available on BIG-IP
 
 
## Example Usage
```hcl

data "bigip_ltm_monitor" "Monitor-TC1" {
  name      = "test-monitor"
  partition = "Common"
}

```      

## Argument Reference

* `name` - (Required) Name of the ltm monitor

* `partition` - (Required) partition of the ltm monitor


## Attributes Reference

Additionally, the following attributes are exported:

* `destination` - id will be full path name of ltm monitor.

* `interval` - Specifies, in seconds, the frequency at which the system issues the monitor check when either the resource is down or the status of the resource is unknown.

* `ip_dscp` - Displays the differentiated services code point (DSCP). DSCP is a 6-bit value in the Differentiated Services (DS) field of the IP header.

* `manual_resume` - Displays whether the system automatically changes the status of a resource to Enabled at the next successful monitor check.

* `reverse` - Instructs the system to mark the target resource down when the test is successful.

* `transparent` - Displays whether the monitor operates in transparent mode.

* `adaptive_limit` - Displays whether adaptive response time monitoring is enabled for this monitor.

* `adaptive` - Displays whether adaptive response time monitoring is enabled for this monitor.
