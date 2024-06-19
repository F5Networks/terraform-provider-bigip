---
page_title: "User guide for deploying per-app AS3 declaration"
description: |-
  User guide for deploying per-app AS3 declaration using **bigip_as3** resource
---

# User Guide for deploying per-app AS3 declaration

~>  **NOTE**  **bigip_as3**  resource supports **Per-Application** mode of AS3 deployment from provider version > `v1.22.1`, more information on **Per-Application** mode can be found [Per-App](https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/per-app-declarations.html)

~> **NOTE** For Supporting AS3 Per-App mode of deployment, AS3 version on **BIG-IP** should be > **v3.50**

~> **NOTE** For Deploying AS3 JSON in Per-App mode, resource should provide with a attribute [tenant_name](#tenant_name) to be passed for deploying **application/applications** on specified tenant, else random tenant name will be generated.

~> **NOTE** **PerApplication** needs to be turned `true` as a Prerequisite on the Big-IP (BIG-IP AS3 version >3.50) device. For details : <https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/per-app-declarations.html>


## Per-Application Deployment - Example Usage for json file with tenant name

```hcl
resource "bigip_as3" "as3-example1" {
  as3_json    = file("perApplication_example.json")
  tenant_name = "Test"
}
```

Example AS3 declaration for the PerApp mode of deployment

`perApplication_example.json` 

```json
{
    "Application1": {
        "class": "Application",
        "service": {
            "class": "Service_HTTP",
            "virtualAddresses": [
                "192.0.2.1"
            ],
            "pool": "pool"
        },
        "pool": {
            "class": "Pool",
            "members": [
                {
                    "servicePort": 80,
                    "serverAddresses": [
                        "192.0.2.10",
                        "192.0.2.20"
                    ]
                }
            ]
        }
    },  
    "Application2": {
        "class": "Application",
        "service": {
            "class": "Service_HTTP",
            "virtualAddresses": [
                "192.0.3.2"
            ],
            "pool": "pool"
        },
        "pool": {
            "class": "Pool",
            "members": [
                {
                    "servicePort": 80,
                    "serverAddresses": [
                        "192.0.3.30",
                        "192.0.3.40"
                    ]
                }
            ]
        }
    }
}
```