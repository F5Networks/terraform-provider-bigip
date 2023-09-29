---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_cipher_rule"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_cipher_rule resource
---

# bigip\_ltm\_cipher\_rule
`bigip_ltm_cipher_rule` Manages F5 BIG-IP LTM cipher rule using iControl REST.

## Example Usage

```hcl
resource "bigip_ltm_cipher_rule" "test_cipher_rule" {
  name                 = "/Common/test_cipher_rule"
  cipher               = "TLS13-AES128-GCM-SHA256:TLS13-AES256-GCM-SHA384"
  dh_groups            = "P256:P384:FFDHE2048:FFDHE3072:FFDHE4096"
  signature_algorithms = "DEFAULT"
}
```

## Argument Reference

* `name` - (Required,type `string`) Name of the Cipher Rule. Name should be in pattern `partition` + `cipher_rule_name`

* `description` - (Optional,type `string`) The Partition in which the Cipher Rule will be created.

* `cipher` - (Required,type `string`) Specifies one or more Cipher Suites used,this is a colon (:) separated string of cipher suites. example, `TLS13-AES128-GCM-SHA256:TLS13-AES256-GCM-SHA384`.

* `dh_groups` - (Optional,type `string`) Specifies the DH Groups algorithms, separated by colons (:).

* `signature_algorithms` - (Optional,type `string`) Specifies the Signature Algorithms, separated by colons (:).

## Importing
An existing cipher rule can be imported into this resource by supplying the cipher rule full path name  ex : `/partition/name`
An example is below:
```sh
$ terraform import bigip_ltm_cipher_rule.test_cipher_rule /Common/test_cipher_rule
```