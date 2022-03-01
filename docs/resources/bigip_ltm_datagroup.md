---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_datagroup"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_datagroup resource
---

# bigip\_ltm\_datagroup

`bigip_ltm_datagroup` Manages internal (in-line) datagroup configuration

Resource should be named with their`full path`. The full path is the combination of the `partition + name` of the resource, for example `/Common/my-datagroup`.


### Example Usage for creating Internal datagroup.

```hcl
resource "bigip_ltm_datagroup" "datagroup" {
  name = "/Common/dgx2"
  type = "string"
  record {
    name = "abc.com"
    data = "pool1"
  }
  record {
    name = "test"
    data = "123"
  }
}
```

### Example Usage for creating External datagroup.

```hcl

resource "bigip_ltm_datagroup" "DG-TC2" {
  name        = "/Common/dgtc2"
  type        = "string"
  internal    = false
  records_src = "/pathtodatagroup/ext_dg_string.txt"
}
```      

## Argument Reference

* `name` - (Required) Name of the datagroup

* `type` - (Required) datagroup type (applies to the `name` field of the record), supports: `string`, `ip` or `integer`

* `records_src` - (Optional, `string`) Path to a file with records in it,The file should be well-formed,it includes records, one per line,that resemble the following format "key separator value". For example, `foo := bar`.
This should be used in conjunction with `internal` attribute set `false`

* `internal` - (Optional,`bool`) Set `false` if you want to Create External Datagroups. default is `true`,means creates internal datagroup.

* `record` - (Optional) a set of `name` and `data` attributes, name must be of type specified by the `type` attributed (`string`, `ip` and `integer`), data is optional and can take any value, multiple `record` sets can be specified as needed.

  * `name` - (Required if `record` defined), sets the value of the record's `name` attribute, must be of type defined in `type` attribute

  * `data` - (Optional if `record` defined), sets the value of the record's `data` attribute, specifying a value here will create a record in the form of `name := data`
