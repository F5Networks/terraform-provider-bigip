---
layout: "bigip"
page_title: "BIG-IP: bigip_gtm_datacenter"
subcategory: "Global Traffic Manager(GTM)"
description: |-
  Provides details about bigip_gtm_datacenter data source
---

# bigip\_gtm\_datacenter

Use this data source (`bigip_gtm_datacenter`) to look up an existing GTM datacenter on the BIG-IP. This is useful when datacenters are shared across multiple Terraform workspaces and you need to reference them without managing them as resources.

## Example Usage

```hcl
data "bigip_gtm_datacenter" "dc" {
  name      = "my-datacenter"
  partition = "Common"
}

output "datacenter_enabled" {
  value = data.bigip_gtm_datacenter.dc.enabled
}

output "datacenter_location" {
  value = data.bigip_gtm_datacenter.dc.location
}
```

## Argument Reference

* `name` - (Required) Name of the GTM datacenter.
* `partition` - (Required) Partition of the GTM datacenter.

## Attributes Reference

Additionally, the following attributes are exported:

* `description` - Description of the datacenter.

* `contact` - Contact information for the datacenter.

* `enabled` - Whether the datacenter is enabled.

* `location` - Location of the datacenter.

* `prober_fallback` - Type of prober to use for fallback.

* `prober_preference` - Type of prober to prefer.