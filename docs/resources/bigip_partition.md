---
layout: "bigip"
page_title: "BIG-IP: bigip_partition"
sidebar_current: "docs-bigip-resource-device-x"
description: |-
  Provides details about bigip_partition resource
---

# bigip_partition

`bigip_partition` Manages F5 BIG-IP partitions

## Example Usage

```hcl
resource "bigip_partition" "test-partition" {
  name            = "test-partition"
  description     = "created by terraform"
  route_domain_id = 2
}
```

## Argument Reference

* `name` - (Required,type `string`) Name of the partition.

* `description` - (Optional,type `string`) Description of the partition.

* `route_domain_id` - (Optional,type `number`) Route domain id of the partition.

## Importing

An existing cipher group can be imported into this resource by supplying the cipher rule full path name ex : `/partition/name`

An example is below:

```sh
$ terraform import bigip_partition.test_partition test_partition

```
