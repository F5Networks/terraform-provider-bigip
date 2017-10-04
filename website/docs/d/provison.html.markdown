---
layout: "bigip"
page_title: "BIG-IP: bigip_provision"
sidebar_current: "docs-bigip-datasource-provision-x"
description: |-
    Provides details about bigip  provision resource for BIG-IP
---

# bigip\_provision

`bigip_provision` provides details bout how to enable "ilx", "asm" "apm" resource on BIG-IP
## Elxample Usage


```hcl
provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_provision" "provision-ilx" {
  name = "/Common/ilx"
  fullPath  = "ilx"
  cpuRatio = 0
  diskRatio = 0
  level = "nominal"
  memoryRatio = 0
}
``` 

## Argument Reference

* `bigip_provsion` - Common is the partition and ilx is the module being enabled it could be asm, afm apm etc.
* `Common/ilx` - Common is the partition and ilx is the module being enabled it could be asm, afm apm etc.
* `cpuRatio` - m, afm apm etc.


      
