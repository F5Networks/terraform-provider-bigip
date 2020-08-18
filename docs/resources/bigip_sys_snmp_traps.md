---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_snmp_traps"
sidebar_current: "docs-bigip-resource-snmp_traps-x"
description: |-
    Provides details about bigip  snmp_traps resource for BIG-IP
---

# bigip\_sys\_snmp\_traps

`bigip_sys_snmp_traps` provides details bout how to enable snmp_traps resource on BIG-IP
## Example Usage


```hcl
resource "bigip_sys_snmp_traps" "snmp_traps" {
  name        = "snmptraps"
  community   = "f5community"
  host        = "195.10.10.1"
  description = "Setup snmp traps"
  port        = 111
}

```

## Argument Reference

* `name` -  (Optional) Name of the snmp trap.

* `community` - (Optional) Specifies the community string used for this trap.

* `host` - The host the trap will be sent to.

* `description` - (Optional) The port that the trap will be sent to.

* `port` - (Optional) User defined description.
