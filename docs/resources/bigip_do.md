---
layout: "bigip"
page_title: "BIG-IP: bigip_do"
sidebar_current: "docs-bigip-resource-x"
description: |-
    Provides details about bigip do resource
---

# bigip_do

`bigip_do` provides details about bigip do resource

This resource is helpful to configure do declarative JSON on BIG-IP.
## Example Usage


```hcl

resource "bigip_do"  "do-example" {
     do_json = "${file("example.json")}"
     timeout = 15
 }

```

## Argument Reference


* `do_json` - (Required) Name of the of the Declarative DO JSON file

* `timeout(minutes)` - (optional) timeout to keep polling DO endpoint until Bigip is provisioned by DO.( Default timeout is 20 minutes ) 

*  Note: Delete method is not supported by DO, so terraform destroy won't delete configuration in bigip but we will set the terrform
   state to empty and won't throw error.

* `example.json` - Example of DO Declarative JSON

```hcl
{
    "schemaVersion": "1.0.0",
    "class": "Device",
    "async": true,  
    "label": "my BIG-IP declaration for declarative onboarding",
    "Common": {
        "class": "Tenant",
        "hostname": "bigip.example.com",
        "myLicense": {
            "class": "License",
            "licenseType": "regKey",
            "regKey": "xxxx"
        }, 
        "admin": {
            "class": "User",
            "userType": "regular",
            "password": "xxxx",
            "shell": "bash"
        },
        "myProvisioning": {
            "class": "Provision",
            "ltm": "nominal",
            "gtm": "minimum"
        },
        "external": {
            "class": "VLAN",
            "tag": 4093,
            "mtu": 1500,
            "interfaces": [
                {
                    "name": "1.1",
                    "tagged": true
                }
            ],
            "cmpHash": "dst-ip"
        },
        "external-self": {
            "class": "SelfIp",
            "address": "x.x.x.x",
            "vlan": "external",
            "allowService": "default",
            "trafficGroup": "traffic-group-local-only"
        }
    
    }
}
```
* `DO documentation` - https://clouddocs.f5.com/products/extensions/f5-declarative-onboarding/latest/composing-a-declaration.html#sample-declaration-for-a-standalone-big-ip
