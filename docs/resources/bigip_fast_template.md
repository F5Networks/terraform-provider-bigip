---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_template"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_fast_template resource
---

# bigip_fast_template

`bigip_fast_template` This resource will import and create FAST template sets on BIG-IP LTM.
Template set can be imported from zip archive files on the local disk.


## Example Usage

```hcl
resource "bigip_fast_template" "foo-template" {
  name     = "foo_template"
  source   = "foo_template.zip"
  md5_hash = filemd5("foo_template.zip")
}
```      

## Argument Reference

* `name`- (Optional) Name of the FAST template set to be created on to BIGIP

* `source` - (Required) Path to the zip archive file containing FAST template set on Local Disk

* `md5_hash` - (Required) MD5 hash of the zip archive file containing FAST template
