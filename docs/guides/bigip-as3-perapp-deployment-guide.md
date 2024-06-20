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


## Per-Application Deployment Create - Example Usage for json file with tenant name

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
    "schemaVersion": "3.50.0",
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

## Per-Application Deployment Update - 

## Update deployed Applications 

```hcl
resource "bigip_as3" "as3-example1" {
  as3_json    = file("updated_perApplication_example.json")
  tenant_name = "Test"
}
```

Example AS3 declaration for the PerApp mode of deployment

`updated_perApplication_example.json`

```json
{   
    "schemaVersion": "3.50.0",
    "Application1": {
        "class": "Application",
        "service": {
            "class": "Service_HTTP",
            "virtualAddresses": [
                "192.0.2.11"
            ],
            "pool": "pool"
        },
        "pool": {
            "class": "Pool",
            "members": [
                {
                    "servicePort": 80,
                    "serverAddresses": [
                        "192.0.2.11",
                        "192.0.2.21"
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
                "192.0.3.22"
            ],
            "pool": "pool"
        },
        "pool": {
            "class": "Pool",
            "members": [
                {
                    "servicePort": 80,
                    "serverAddresses": [
                        "192.0.3.32",
                        "192.0.3.42"
                    ]
                }
            ]
        }
    }
}
```


## Add new Application
either the new Application can be added in the same AS3 declaration (to make it include all the applications) or new resource can be added to passing only the new application AS3 declaration

```hcl
resource "bigip_as3" "as3-example1" {
  as3_json    = file("perApplication_example.json")
  tenant_name = "Test"
}
```

AS3 declaration including all the Apps for PerApp mode of deployment

`perApplication_example.json` 

```json
{
    "schemaVersion": "3.50.0",
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
    },
    "Application3": {
        "class": "Application",
        "service": {
            "class": "Service_HTTP",
            "virtualAddresses": [
                "192.0.3.3"
            ],
            "pool": "pool"
        },
        "pool": {
            "class": "Pool",
            "members": [
                {
                    "servicePort": 80,
                    "serverAddresses": [
                        "192.0.3.50",
                        "192.0.3.60"
                    ]
                }
            ]
        }
    }
}
```

```hcl
resource "bigip_as3" "as3-example1" {
  as3_json    = file("perApplication_example.json")
  tenant_name = "Test"
}

resource "bigip_as3" "as3-example2" {
  as3_json    = file("perApplication_example2.json")
  tenant_name = "Test"
}
```
AS3 declaration including only the new App for PerApp mode of deployment 

`perApplication_example.json` 

```json
{
    "schemaVersion": "3.50.0",
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


`perApplication_example2.json`

```json
{
    "schemaVersion": "3.50.0",
    "Application3": {
        "class": "Application",
        "service": {
            "class": "Service_HTTP",
            "virtualAddresses": [
                "192.0.3.3"
            ],
            "pool": "pool"
        },
        "pool": {
            "class": "Pool",
            "members": [
                {
                    "servicePort": 80,
                    "serverAddresses": [
                        "192.0.3.50",
                        "192.0.3.60"
                    ]
                }
            ]
        }
    }
}
```


# Delete Existing App

Assuming we have 3 Applications deployed - Application1 and Application2 (via as3-example1 resource) and Application3 via as3-example2 resource.
In order to delete specific App - Application2, we can update the AS3 blob passed in as3-example1


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
    "schemaVersion": "3.50.0",
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
    }
}
```

In order to delete all applications deployed via a separate resources eg:
- In order to delete Application1,Application2 , we can directly run command `terraform destroy -target=bigip_as3.as3-example1`
- In order to delete Application3, we can directly run command `terraform destroy -target=bigip_as3.as3-example2`

If only one App is deployed in a Tenant - on deletion of that App, Tenant will also be deleted.
