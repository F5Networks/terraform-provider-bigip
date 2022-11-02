---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_fastl4"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_fastl4 resource
---

# bigip\_ltm\_profile_fastl4

`bigip_ltm_profile_fastl4` Configures a custom LTM fastL4 profile for use by health checks.

Resources should be named with their `full path`. The full path is the combination of the `partition + name` of the resource (For example `/Common/my-fastl4profile`) or  `partition + directory + name` of the resource  (example: `/Common/test/my-fastl4profile`)

## Example Usage


```hcl
resource "bigip_ltm_profile_fastl4" "profile_fastl4" {
  name                   = "/Common/sjfastl4profile"
  defaults_from          = "/Common/fastL4"
  client_timeout         = 40
  explicitflow_migration = "enabled"
  hardware_syncookie     = "enabled"
  idle_timeout           = "200"
  iptos_toclient         = "pass-through"
  iptos_toserver         = "pass-through"
  keepalive_interval     = "disabled" //This cannot take enabled
}

```      

## Argument Reference

* `name` (Required,type `string`) Name of the LTM fastL4 Profile.The full path is the combination of the `partition + name` of the resource (For example `/Common/my-fastl4profile`) or  `partition + directory + name` of the resource  (example: `/Common/test/my-fastl4profile`)

* `defaults_from` - (Optional,type `string`) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `late_binding` - (Optional,type `string`) Enables intelligent selection of a back-end server or pool, using an iRule to make the selection. The default is `disabled`.

* `client_timeout` - (Optional,type `int`) Specifies late binding client timeout in seconds. This setting specifies the number of seconds allowed for a client to transmit enough data to select a server when late binding is enabled. If it expires timeout-recovery mode will dictate what action to take.

* `explicitflow_migration` - (Optional,type `string`)Enables or disables late binding explicit flow migration that allows iRules to control when flows move from software to hardware. Explicit flow migration is disabled by default hence BIG-IP automatically migrates flows from software to hardware.

* `hardware_syncookie` - (Optional,type `string`) Enables or disables hardware SYN cookie support when PVA10 is present on the system. Note that when you set the hardware syncookie option to enabled, you may also want to set the following bigdb database variables using the "/sys modify db" command, based on your requirements: pva.SynCookies.Full.ConnectionThreshold (default: 500000), pva.SynCookies.Assist.ConnectionThreshold (default: 500000) pva.SynCookies.ClientWindow (default: 0). The default value is disabled.

* `idle_timeout` - (Optional,type `string`) Specifies an idle timeout in seconds. This setting specifies the number of seconds that a connection is idle before the connection is eligible for deletion.When you specify an idle timeout for the Fast L4 profile, the value must be greater than the bigdb database variable Pva.Scrub time in msec for it to work properly.The default value is 300 seconds.

* `iptos_toclient` - (Optional,type `string`) Specifies an IP ToS number for the client side. This option specifies the Type of Service level that the traffic management system assigns to IP packets when sending them to clients. The default value is 65535 (pass-through), which indicates, do not modify.

* `iptos_toserver`  - (Optional,type `string`) Specifies an IP ToS number for the server side. This setting specifies the Type of Service level that the traffic management system assigns to IP packets when sending them to servers. The default value is 65535 (pass-through), which indicates, do not modify.

* `keepalive_interval` - (Optional,type `string`) Specifies the keep alive probe interval, in seconds. The default value is disabled (0 seconds).

* `tcp_handshake_timeout` - (Optional,type `string`) Specifies the acceptable duration for a TCP handshake, that is, the maximum idle time between a client synchronization (SYN) and a client acknowledgment (ACK).The default is `5 seconds`.

* `loose_initiation` - (Optional,type `string`) Specifies, when checked (enabled), that the system initializes a connection when it receives any TCP packet, rather that requiring a SYN packet for connection initiation. The default is disabled. We recommend that if you enable the Loose Initiation option, you also enable the Loose Close option.

* `loose_close` - (Optional,type `string`) Specifies, when checked (enabled), that the system closes a loosely-initiated connection when the system receives the first FIN packet from either the client or the server. The default is disabled.

* `receive_windowsize` - (Optional,type `int`) Specifies the amount of data the BIG-IP system can accept without acknowledging the server. The default is 0 (zero).

## Import

BIG-IP LTM fastl4 profiles can be imported using the `name`, e.g.

```
$ terraform import bigip_ltm_profile_fastl4.test-fastl4 /Common/test-fastl4
```
