---
layout: "bigip"
page_title: "BIG-IP: bigip_ssl_key"
subcategory: "System"
description: |-
  Provides details about bigip_ssl_key resource
---

# bigip_ssl_key

`bigip_ssl_key` This resource will import SSL certificate key on BIG-IP LTM. 
Certificate key can be imported from certificate key files on the local disk, in PEM format


## Example Usage


```hcl

resource "bigip_ssl_key" "test-key" {
  name      = "serverkey.key"
  content   = file("serverkey.key")
  partition = "Common"
}

```      

## Argument Reference


* `name`- (Required,type `string`) Name of the SSL Certificate key to be Imported on to BIGIP

* `content` - (Required) Content of certificate key on Local Disk,path of SSL certificate key will be provided to terraform `file` function 

* `partition` - (Optional,type `string`) Partition on to SSL Certificate key to be imported. The parameter is not required when running terraform import operation. In such case the name must be provided in `full_path` format.

