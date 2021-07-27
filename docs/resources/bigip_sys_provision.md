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
* `Common/ilx` - Common is the partition and ilx is the module being enabled it could be asm, afm apm etc.
* `cpuRatio` - how much cpu resources you need for this resource
* `diskRatio` - how much disk space you want to allocate for this resource.
* `memoryRatio` - how much memory you want to deidcate for this resource
