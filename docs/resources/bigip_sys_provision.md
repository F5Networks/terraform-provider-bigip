---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_provision"
sidebar_current: "docs-bigip-resource-provision-x"
description: |-
   Provides details about bigip  provision resource for BIG-IP
---

# bigip\_sys\_provision

`bigip_sys_provision` provides details bout how to enable "ilx", "asm" "apm" resource on BIG-IP
## Example Usage


```hcl
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}
resource "bigip_sys_provision" "test-provision" {
  name         = "TEST_ASM_PROVISION_NAME"
  full_path    = "asm"
  cpu_ratio    = 0
  disk_ratio   = 0
  level        = "none"
  memory_ratio = 0
}

```

## Argument Reference

* `bigip_sys_provision` - Is the resource which is used to provision big-ip modules like asm, afm, ilx etc
* `full_path` - Specifies the module being enabled. It can be one among the list [afm, am, apm, asm, avr, cgnat, fps, gtm, ilx, lc, ltm, pem, swg, urldb, sslo, vcmp]
* `cpuRatio` - (Optional) Use this option only when the level option is set to custom. F5 recommends that you do not modify this option. The default value is zero.
* `diskRatio` - (Optional) Use this option only when the level option is set to custom. F5 recommends that you do not modify this option. The default value is zero.
* `memoryRatio` - (Optional) Use this option only when the level option is set to custom. F5 recommends that you do not modify this option. The default value is zero.
* `level` - (Optional) Specifies the level of resources that you want to provision for a module. The list of allowed values are [minimum, nominal, dedicated, custom, none]
