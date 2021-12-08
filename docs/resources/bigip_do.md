---
layout: "bigip"
page_title: "BIG-IP: bigip_do"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_do resource
---

# bigip_do

`bigip_do` provides details about bigip do resource

This resource is helpful to configure do declarative JSON on BIG-IP.
## Example Usage


```hcl

resource "bigip_do" "do-example" {
  do_json = "${file("example.json")}"
  timeout = 15
}

```

## Argument Reference


* `do_json` - (Required) Name of the of the Declarative DO JSON file
 
* `bigip_address` - (optional) IP Address of BIGIP Host to be used for this resource,this is optional parameter.
whenever we specify this parameter it gets overwrite provider configuration

* `bigip_user` - (optional) UserName of BIGIP host to be used for this resource,this is optional parameter.
whenever we specify this parameter it gets overwrite provider configuration

* `bigip_port` - (optional) Port number of BIGIP host to be used for this resource,this is optional parameter.
whenever we specify this parameter it gets overwrite provider configuration

* `bigip_password` - (optional) Password of  BIGIP host to be used for this resource,this is optional parameter.
whenever we specify this parameter it gets overwrite provider configuration

* `timeout(minutes)` - (optional) timeout to keep polling DO endpoint until Bigip is provisioned by DO.( Default timeout is 20 minutes )

~> **Note:** If we want to replace provider BIGIP with other BIGIPs details we can specify with `bigip_address`,
`bigip_user`,`bigip_port` and `bigip_password`. All Must be specified in such scenario.
   
~> **Note:** Delete method is not supported by DO, so terraform destroy won't delete configuration in bigip but we will set the terrform
   state to empty and won't throw error.


* `example.json` - Example of DO Declarative JSON

```json
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
