---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_fastl4"
sidebar_current: "docs-bigip-resource-profile_fastl4-x"
description: |-
    Provides details about bigip_ltm_profile_fastl4 resource
---

# bigip\_ltm\_profile_fastl4

`bigip_ltm_profile_fastl4` Configures a custom profile_fastl4 for use by health checks.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage


```hcl
resource "bigip_ltm_profile_fastl4" "profile_fastl4" {
  name                   = "/Common/sjfastl4profile"
  partition              = "Common"
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

* `name` (Required) Name of the profile_fastl4

* `defaults_from` - (Optional) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `partition` - (Optional) Displays the administrative partition within which this profile resides

* `client_timeout` - (Optional) Specifies late binding client timeout in seconds. This setting specifies the number of seconds allowed for a client to transmit enough data to select a server when late binding is enabled. If it expires timeout-recovery mode will dictate what action to take.

* `explicitflow_migration` - (Optional) Enables or disables late binding explicit flow migration that allows iRules to control when flows move from software to hardware. Explicit flow migration is disabled by default hence BIG-IP automatically migrates flows from software to hardware.

* `hardware_syncookie` - (Optional) Enables or disables hardware SYN cookie support when PVA10 is present on the system. Note that when you set the hardware syncookie option to enabled, you may also want to set the following bigdb database variables using the "/sys modify db" command, based on your requirements: pva.SynCookies.Full.ConnectionThreshold (default: 500000), pva.SynCookies.Assist.ConnectionThreshold (default: 500000) pva.SynCookies.ClientWindow (default: 0). The default value is disabled.

* `idle_timeout` - (Optional) Specifies an idle timeout in seconds. This setting specifies the number of seconds that a connection is idle before the connection is eligible for deletion.When you specify an idle timeout for the Fast L4 profile, the value must be greater than the bigdb database variable Pva.Scrub time in msec for it to work properly.The default value is 300 seconds.

* `iptos_toclient` - (Optional) Specifies an IP ToS number for the client side. This option specifies the Type of Service level that the traffic management system assigns to IP packets when sending them to clients. The default value is 65535 (pass-through), which indicates, do not modify.

* `iptos_toserver`  - (Optional) Specifies an IP ToS number for the server side. This setting specifies the Type of Service level that the traffic management system assigns to IP packets when sending them to servers. The default value is 65535 (pass-through), which indicates, do not modify.

* `keepalive_interval` - (Optional) Specifies the keep alive probe interval, in seconds. The default value is disabled (0 seconds).
