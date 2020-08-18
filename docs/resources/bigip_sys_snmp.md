---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_snmp"
sidebar_current: "docs-bigip-resource-snmp-x"
description: |-
    Provides details about bigip  snmp resource for BIG-IP
---

# bigip\_sys\_snmp

`bigip_sys_snmp` provides details bout how to enable "ilx", "asm" "apm" resource on BIG-IP
## Example Usage


```hcl

resource "bigip_sys_snmp" "snmp" {
  sys_contact      = " NetOPsAdmin s.shitole@f5.com"
  sys_location     = "SeattleHQ"
  allowedaddresses = ["202.10.10.2"]
}

```

## Argument Reference

* `sys_contact` -  (Optional) Specifies the contact information for the system administrator.

* `sys_location` - Describes the system's physical location.

* `allowedaddresses` - Configures hosts or networks from which snmpd can accept traffic. Entries go directly into hosts.allow.
