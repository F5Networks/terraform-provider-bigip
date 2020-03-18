---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_tcp"
sidebar_current: "docs-bigip-resource-profile_tcp-x"
description: |-
    Provides details about bigip_ltm_profile_tcp resource
---

# bigip\_ltm\_profile_tcp

`bigip_ltm_profile_tcp` Configures a custom profile_tcp for use by health checks.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage


```hcl
resource "bigip_ltm_profile_tcp" "sanjose-tcp-lan-profile" {
  name               = "sanjose-tcp-lan-profile"
  idle_timeout       = 200
  close_wait_timeout = 5
  finwait_2timeout   = 5
  finwait_timeout    = 300
  keepalive_interval = 1700
  deferred_accept    = "enabled"
  fast_open          = "enabled"
}

```      

## Argument Reference

* `name` (Required) Name of the profile_tcp

* `partition` - (Optional) Displays the administrative partition within which this profile resides

* `defaults_from` - (Optional) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.


* `idle_timeout` - (Optional) Specifies the number of seconds that a connection is idle before the connection is eligible for deletion. The default value is 300 seconds.

* `close_wait_timeout` - (Optional) Specifies the number of seconds that a connection remains in a LAST-ACK state before quitting. A value of 0 represents a term of forever (or until the maxrtx of the FIN state). The default value is 5 seconds.

* `finwait_timeout` - (Optional) Specifies the number of seconds that a connection is in the FIN-WAIT-1 or closing state before quitting. The default value is 5 seconds. A value of 0 (zero) represents a term of forever (or until the maxrtx of the FIN state). You can also specify immediate or indefinite.

* `finwait_2timeout` - (Optional) Specifies the number of seconds that a connection is in the FIN-WAIT-2 state before quitting. The default value is 300 seconds. A value of 0 (zero) represents a term of forever (or until the maxrtx of the FIN state).

* `keepalive_interval` - (Optional) Specifies the keep alive probe interval, in seconds. The default value is 1800 seconds.


* `fast_open` - (Optional) When enabled, permits TCP Fast Open, allowing properly equipped TCP clients to send data with the SYN packet.

* `deferred_accept` - (Optional) Specifies, when enabled, that the system defers allocation of the connection chain context until the client response is received. This option is useful for dealing with 3-way handshake DOS attacks. The default value is disabled.
