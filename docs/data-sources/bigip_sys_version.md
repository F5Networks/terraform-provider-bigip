---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_version"
subcategory: "System"
description: |-
  Provides details about bigip_sys_version data source
---

# bigip\_sys\_version

Use this data source (`bigip_sys_version`) to retrieve the BIG-IP system's running TMOS version and platform information. `sys/version` (and the platform information supplemented from `sys/hardware`) is read-only system information -- there's nothing to filter by -- so this data source takes no arguments.

This is useful for migration tooling that needs to determine the currently running version/build in order to select a matching r-Series/VELOS tenant image (see [`scripts/inventory-tmos-version.sh`](https://github.com/F5Networks/terraform-provider-bigip/blob/main/scripts/inventory-tmos-version.sh), which derives the same recommended tenant image filename via a `bigip_command`-based `tmsh` parse instead).

## Example Usage

```hcl
data "bigip_sys_version" "current" {}

output "tmos_version" {
  value = data.bigip_sys_version.current.version
}

output "tenant_image_hint" {
  value = "BIGIP-${data.bigip_sys_version.current.version}-${data.bigip_sys_version.current.build}.<TYPE>-F5OS"
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `version` - Running TMOS version, e.g. `17.5.1`.

* `build` - Build number of the running TMOS version, e.g. `0.0.7`.

* `product` - Product name, e.g. `BIG-IP`.

* `edition` - Release edition, e.g. `Final`, `Hotfix`, `Engineering Hotfix`.

* `platform` - Platform marketing name (from `sys/hardware`), e.g. `BIG-IP Virtual Edition` or a specific appliance model. Empty if `sys/hardware`'s platform information could not be retrieved.
