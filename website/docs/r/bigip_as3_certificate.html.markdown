---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_certificate"
sidebar_current: "docs-bigip-datasource-certificate-x"
description: |-
   Provides details about bigip_as3_cert datasource
---

# bigip\_as3\_cert

`bigip_as3_cert` Manages a Certificate class, which  contain your servers as well as health monitors and load balancing methods and more.

In the following example, our pool is web_pool, it’s using the default HTTP health monitor, and includes two servers on port 80.

## Example Usage


```hcl
data "bigip_as3_cert" "exmpcert" {
  name = "exmpcert"
  remark = "in practice we recommend using a passphrase"
  certificate = "${file("servercert.crt")}"
  private_key = "${file("serverkey.key")}"
  passphrase {
    ciphertext = "ZjVmNQ=="
    protected = "eyJhbGciOiJkaXIiLCJlbmMiOiJub25lIn0"
  }
}
```

## Argument Reference

* `name` - (Required) Name of the certificate

* `certificate` - X.509 public-key certificate 

* `private_key` - (Optional) Private key matching certificate’s public key

* `label` - (Optional) Optional friendly name for this object

* `remark` - (Optional) Arbitrary (brief) text pertaining to this object

* Below attributes needs to be configured under passphrase option

* `ciphertext`- (Optional) Put base64url(data_value) here

* `protected` - (Optional) JOSE header: alg=dir, enc=(none|f5sv); default enc=none (encoded default is ‘protected’=’eyJhbGciOiJkaXIiLCJlbmMiOiJub25lIn0’, use with secret simply base64url-encoded into ‘ciphertext’). If you see ‘protected’=’eyJhbGciOiJkaXIiLCJlbmMiOiJmNXN2In0’, ‘ciphertext’ contains base64url-encoded SecureVault cryptogram 
