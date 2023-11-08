---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_cipher_group"
subcategory: "Local Traffic Manager(LTM)"
description: |-
Provides details about bigip_ltm_cipher_group resource
---

# bigip\_ltm\_cipher\_group

`bigip_ltm_cipher_group` Manages F5 BIG-IP LTM cipher group using iControl REST.

## Example Usage

```hcl
resource "bigip_ltm_cipher_group" "test-cipher-group" {
  name     = "/Common/test-cipher-group-01"
  allow    = ["/Common/f5-aes"]
  require  = ["/Common/f5-quic"]
  ordering = "speed"
}
```

## Argument Reference

* `name` - (Required,type `string`) Name of the Cipher group. Name should be in pattern `partition` + `cipher_group_name`

* `allow` - (Optional,type `list` of `strings` ) Specifies the configuration of the allowed groups of ciphers. You can select a cipher rule from the Available Cipher Rules list.

* `require` - (Optional,type `list` of `string`) Specifies the configuration of the restrict groups of ciphers. You can select a cipher rule from the Available Cipher Rules list.

* `ordering` - (Optional,type `string`) Controls the order of the Cipher String list in the Cipher Audit section. Options are Default, Speed, Strength, FIPS, and Hardware. The rules are processed in the order listed.

## Importing
An existing cipher group can be imported into this resource by supplying the cipher rule full path name ex : `/partition/name`
An example is below:
```sh
$ terraform import bigip_ltm_cipher_group.test_cipher_group /Common/test_cipher_group

```