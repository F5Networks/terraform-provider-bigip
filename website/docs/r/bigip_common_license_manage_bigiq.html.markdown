---
layout: "bigip"
page_title: "BIG-IP: bigip_common_license_manage_bigiq"
sidebar_current: "docs-bigip-resource-node-x"
description: |-
  Provides details about bigip_common_license_manage_bigiq resource
---

# bigip_common_license_manage_bigiq


`bigip_common_license_manage_bigiq` This Resource is used for BIGIP/Provider License Management from BIGIQ


## Example Usage


```hcl

# MANAGED Regkey Pool
resource "bigip_common_license_manage_bigiq" "test_example" {
  bigiq_address = var.bigiq
  bigiq_user = var.bigiq_un
  bigiq_password = var.bigiq_pw
  license_poolname = "regkeypool_name"
  assignment_type = "MANAGED"
}

# UNMANAGED Regkey Pool
resource "bigip_common_license_manage_bigiq" "test_example" {
  bigiq_address = var.bigiq
  bigiq_user = var.bigiq_un
  bigiq_password = var.bigiq_pw
  license_poolname = "regkeypool_name"
  assignment_type = "UNMANAGED"
} 

# UNMANAGED Utility Pool
resource "bigip_common_license_manage_bigiq" "test_example" {
  bigiq_address = var.bigiq
  bigiq_user = var.bigiq_un
  bigiq_password = var.bigiq_pw
  license_poolname = "utilitypool_name"
  assignment_type = "UNMANAGED"
  unit_of_measure = "yearly"
  skukeyword1 = "BTHSM200M"
}

# UNREACHABLE Regkey Pool
resource "bigip_common_license_manage_bigiq" "test_example" {
  bigiq_address="xxx.xxx.xxx.xxx"
  bigiq_user="xxxx"
  bigiq_password="xxxxx"
  license_poolname = "regkey_pool_name"
  assignment_type = "UNREACHABLE"
  mac_address = "FA:16:3E:1B:6D:32"
  hypervisor = "azure"
}
```      

## Argument Reference

* `bigiq_address` - (Required) BIGIQ License Manager IP Address, variable type `string`

* `bigiq_user` - (Required) BIGIQ License Manager username, variable type `string`

* `bigiq_password` - (Required) BIGIQ License Manager password.  variable type `string`

* `bigiq_port` - (Optional) type `int`, BIGIQ License Manager Port number, specify if port is other than `443`

* `bigiq_token_auth` - (Optional) type `bool`, if set to `true` enables Token based Authentication,default is `false`

* `bigiq_login_ref` - (Optional) BIGIQ Login reference for token authentication

* `assignment_type` - (Required) The type of assignment, which is determined by whether the BIG-IP is unreachable, unmanaged, or managed by BIG-IQ. Possible values: “UNREACHABLE”, “UNMANAGED”, or “MANAGED”.

* `license_poolname` - (Required) A name given to the license pool. type `string`

* `unit_of_measure` - (Optional, Required for `Utility` licenseType) The units used to measure billing. For example, “hourly” or “daily”. Type `string`

* `skukeyword1` - (Optional) An optional offering name. type `string`

* `skukeyword2` - (Optional) An optional offering name. type `string`

* `mac_address` - (Optional, Required Only for `unreachable BIG-IP`) MAC address of the BIG-IP. type `string` 

* `hypervisor` - (Optional,Required Only for `unreachable BIG-IP`) Identifies the platform running the BIG-IP VE. Possible values: “aws”, “azure”, “gce”, “vmware”, “hyperv”, “kvm”, or “xen”. type `string`

* `tenant` - (Optional) For an unreachable BIG-IP, you can provide an optional description for the assignment in this field.

* `key` - (Optional) License Assignment is done with specified `key`, supported only with RegKeypool type License assignement. type `string`


