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

resource "bigip_as3"  "as3-example" {
     as3_json = "${file("example.json")}"
     tenant_name = "as3"
 }

```

## Argument Reference


* `as3_json` - (Required) Name of the of the Declarative AS3 JSON file

* `tenant_name` - (Required) This is the partition name where the application services will be configured.

* `example.json` - Example of AS3 Declarative JSON

```hcl
{
   "class": "AS3",
   "action": "deploy",
   "persist": true,
   "declaration": {
      "class": "ADC",
      "schemaVersion": "3.0.0",
      "id": "urn:uuid:33045210-3ab8-4636-9b2a-c98d22ab915d",
      "label": "Sample 1",
      "remark": "Simple HTTP application with RR pool",
      "as3": {
         "class": "Tenant",
         "A1": {
            "class": "Application",
            "template": "http",
            "serviceMain": {
               "class": "Service_HTTP",
               "virtualAddresses": [
                  "10.0.1.10"
               ],
               "pool": "web_pool"
            },
            "web_pool": {
               "class": "Pool",
               "monitors": [
                  "http"
               ],
               "members": [{
                  "servicePort": 80,
                  "serverAddresses": [
                     "192.0.1.10",
                     "192.0.1.11"
                  ]
               }]
            }
         }
      }
   }
}
```
* `AS3 documentation` - https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/composing-a-declaration.html
