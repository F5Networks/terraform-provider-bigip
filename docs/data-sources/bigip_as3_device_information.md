---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_device_information"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_as3_device_information data source
---

# bigip\_as3\_device\_information

Use this data source (`bigip_as3_device_information`) to retrieve the AS3 declaration currently applied to a tenant on a BIG-IP device.

## Example Usage

```hcl
data "bigip_as3_device_information" "example" {
  tenant = "Sample_tenant"
}

output "as3_declaration" {
  value = data.bigip_as3_device_information.example.as3_json
}
```

### Filtering by application

```hcl
data "bigip_as3_device_information" "example" {
  tenant       = "Sample_tenant"
  applications = ["app1", "app2"]
}
```

## Argument Reference

* `tenant` - (Required) The specific AS3 tenant to retrieve configuration for. Accepts either `tenant_name` or `/partition/tenant_name`.
* `applications` - (Optional) List of applications to retrieve from the tenant. Leave empty to fetch all applications for the tenant.

## Attributes Reference

* `as3_json` - JSON string representation of the retrieved AS3 declaration.
