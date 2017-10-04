---
layout: "bigip"
page_title: "BIG-IP: bigip_datagroup"
sidebar_current: "docs-bigip-datasource-datagroup-x"
description: |-
    Provides details about bigip  datagroup  resource for BIG-IP which is used with iRule
---

# bigip\_datagroup

`bigip_datagroup` provides details bout how to populate the data group values which are used with iRule.
## Example Usage


```hcl
provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}
resource "bigip_datagroup" "datagroup1" {
  name = "dgx2"
  type = "string"
  records  {
   name = "abc.com"
   data = "pool1"
   }
}
``` 

## Argument Reference

* `bigip_datagroup` - Is the resource which is used to configure datagroup values.
* `dgx2` - This is the name of the data group and it can be anything as long as it starts with alphanumeric
* `string` - When the above value is 'string' then the data group is string type.
* `abc.com` - This is the name in the record 
* `data` - This is the object in BIG-IP, it can be a pool member, VIP etc.



      
