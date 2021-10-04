---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_provision"
subcategory: "System"
description: |-
  Provides details about module provision resource for BIG-IP
---

# bigip\_sys\_provision

`bigip_sys_provision` Manage BIG-IP module provisioning. This resource will only provision at the standard levels of Dedicated, Nominal, and Minimum.

## Example Usage

```hcl

resource "bigip_sys_provision" "gtm" {
  name         = "gtm"
  cpu_ratio    = 0
  disk_ratio   = 0
  level        = "nominal"
  memory_ratio = 0
}

```

## Argument Reference

* `name` - (Required,type `string`) Name of module to provision in BIG-IP. 
possible options: 
    * afm
    * am
    * apm
    * cgnat
    * asm
    * avr
    * dos
    * fps
    * gtm
    * ilx
    * lc
    * ltm
    * pem
    * sslo
    * swg
    * urldb
    
* `level` - (Optional,type `string`) Sets the provisioning level for the requested modules. Changing the level for one module may require modifying the level of another module. For example, changing one module to `dedicated` requires setting all others to `none`. Setting the level of a module to `none` means the module is not activated.
default is `nominal`
possible options: 
    * nominal
    * minimum
    * none
    * dedicated

* `cpu_ratio` - (Optional,type `int`) Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none

* `disk_ratio` - (Optional,type `int`)  Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none

* `memory_ratio` - (Optional,type `int`)  Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none
