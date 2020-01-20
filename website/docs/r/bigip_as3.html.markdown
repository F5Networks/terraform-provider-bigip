---
layout: "bigip"
page_title: "BIG-IP: bigip_as3"
sidebar_current: "docs-bigip-resource-x"
description: |-
    Provides details about bigip as3 resource
---

# bigip_as3

`bigip_as3` provides details about bigip as3 resource

This resource is helpful to configure as3 declarative JSON on BIG-IP.
## Example Usage


```hcl

variable "tenant_name" {
  default = "Sample_01"
}

resource "bigip_as3" "as3-demo1" {
  as3_json = templatefile(
    "as3_example.tmpl",
    {
      tenant_name = jsonencode(var.tenant_name)
    })
  tenant_name = var.tenant_name
}

```

## Argument Reference


* `as3_json` - (Required) Path/Filename of Declarative AS3 JSON template file used with builtin ```templatefile``` function 

* `tenant_name` - (Required) Tenant name used to set the terraform state changes for as3 resource

* `as3_example.tmpl` - Example template file  AS3 Declarative JSON

```hcl

 {
     "class": "AS3",
     "action": "deploy",
     "persist": true,
     "declaration": {
         "class": "ADC",
         "schemaVersion": "3.0.0",
         "id": "example-declaration-01",
         "label": "Sample 1",
         "remark": "Simple HTTP application with round robin pool",
         ${tenant_name}: {
             "class": "Tenant",
             "defaultRouteDomain": 0,
             "Application_1": {
                 "class": "Application",
                 "template": "http",
             "serviceMain": {
                 "class": "Service_HTTP",
                 "virtualAddresses": [
                     "10.0.2.10"
                 ],
                 "pool": "web_pool"
                 },
                 "web_pool": {
                     "class": "Pool",
                     "monitors": [
                         "http"
                     ],
                     "members": [
                         {
                             "servicePort": 80,
                             "serverAddresses": [
                                 "192.0.1.100",
                                 "192.0.1.110"
                             ]
                         }
                     ]
                 }
             }
         }
     }
 }

```
* `AS3 documentation` - https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/composing-a-declaration.html
