---
layout: "bigip"
page_title: "BIG-IP: bigip_as3"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_as3 resource
---

# bigip_as3

`bigip_as3` provides details about bigip as3 resource

This resource is helpful to configure AS3 declarative JSON on BIG-IP.

~> This Resource also supports **Per-Application** mode of AS3 deployment, more information on **Per-Application** mode can be found [Per-App](https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/per-app-declarations.html)

-> For Supporting AS3 Per-App mode of deployment, AS3 version on BIG-IP should be > **v3.50**

~> For Deploying AS3 JSON in Per-App mode, this resource provided with a attribute [tenant_name](#tenant_name) to be passed to add application on specified tenant, else random tenant name will be generated.


As3 Declaration can be deployed in Traditional way as well as Per-Application Way : 

- Traditional Way - Entire Declaration needs to be passed in during the create and update call along with the tenant details in the declaration.
- Per-Application Way - Only application details needs to be passed in the as3_json. Tenant name needs to be passed else random tenant name will be generated. 

**Note:** : PerApplication needs to be turned `true` as a Prerequisite on the Big-IP (BIG-IP AS3 version >3.50) device. For details : <https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/per-app-declarations.html>

## Example Usage 

```hcl

# Example Usage for json file
resource "bigip_as3" "as3-example1" {
  as3_json = file("example1.json")
}

# Example Usage for json file with tenant filter
resource "bigip_as3" "as3-example1" {
  as3_json      = file("example2.json")
  tenant_filter = "Sample_03"
}
```

# Per-Application Deployment - Example Usage for json file with tenant name
resource "bigip_as3" "as3-example1" {
  as3_json   = file("perApplication_example.json")
  tenant_name = "Test"
}

# Per-Application Deployment - Example Usage for json file without tenant name - Tenant with Random name is generated in this case
resource "bigip_as3" "as3-example1" {
  as3_json   = file("perApplication_example.json")
}

# Per-Application Deployment - Delete Example
resource "bigip_as3" "as3-example1" {
  tenant_name = "Test"
  as3_json   = file("as3_per_app_example1.json")

}

resource "bigip_as3" "as3-example2" {
  tenant_name = "Test"
  as3_json   = file("as3_per_app_example2.json")
  
}

On running above 2 resources , we will be able to deploy Applications - `path_app1 , path_app2` on Tenant `Test`
now, if we run `terraform destroy -target=bigip_as3.as3-example2` , only `path_app2` will be deleted from Tenant. 


resource "bigip_as3"  "as3-example1" {
	tenant_name = "dmz"
    as3_json = "${file("` + dir + `/../examples/as3/as3_per_app_example1.json")}"
}

resource "bigip_as3"  "as3-example1" {
	tenant_name = "dmz"
    as3_json = "${file("` + dir + `/../examples/as3/as3_per_app_example3.json")}"
}

On running above calls , first call will create 2 applications `path_app1 , path_app2` and second call , will remove `path_app2` from Tenant

[perApplication as3](#perApplication_example)

```hcl

# Per-Application Deployment - Example Usage for json file with tenant name
resource "bigip_as3" "as3-example1" {
  as3_json    = file("perApplication_example.json")
  tenant_name = "Test"
}

# Per-Application Deployment - Example Usage for json file without tenant name - Random tenant name is generated in this case
resource "bigip_as3" "as3-example1" {
  as3_json = file("perApplication_example.json")
}
```

## Argument Reference

* `as3_json` - (Required) Path/Filename of Declarative AS3 JSON which is a json file used with builtin ```file``` function

* `tenant_filter` - (Optional) If there are multiple tenants on a BIG-IP, this attribute helps the user to set a particular tenant to which he want to reflect the changes. Other tenants will neither be created nor be modified.

* `tenant_name` - (Optional) Name of Tenant. This name is used only in the case of Per-Application Deployment. If it is not provided, then a random name would be generated.

* `per_app_mode` - (Computed) - Will specify whether is deployment is done via Per-Application Way or Traditional Way

* `tenant_list` - (Optional) - List of tenants currently deployed on the Big-Ip

* `application_list` - (Optional) - List of applications currently deployed on the Big-Ip

* `ignore_metadata` - (Optional) Set True if you want to ignore metadata changes during update. By default it is set to false

* `as3_example1.json` - Example  AS3 Declarative JSON file with single tenant

```json

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
         "Sample_01": {
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
* `as3_example2.json` - Example  AS3 Declarative JSON file with multiple tenants

```json

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
         "Sample_02": {
             "class": "Tenant",
             "defaultRouteDomain": 0,
             "Application_2": {
                 "class": "Application",
                 "template": "http",
             "serviceMain": {
                 "class": "Service_HTTP",
                 "virtualAddresses": [
                     "10.2.2.10"
                 ],
                 "pool": "web_pool2"
                 },
                 "web_pool2": {
                     "class": "Pool",
                     "monitors": [
                         "http"
                     ],
                     "members": [
                         {
                             "servicePort": 80,
                             "serverAddresses": [
                                 "192.2.1.100",
                                 "192.2.1.110"
                             ]
                         }
                     ]
                 }
             }
         },
         "Sample_03": {
             "class": "Tenant",
             "defaultRouteDomain": 0,
             "Application_3": {
                 "class": "Application",
                 "template": "http",
             "serviceMain": {
                 "class": "Service_HTTP",
                 "virtualAddresses": [
                     "10.1.2.10"
                 ],
                 "pool": "web_pool3"
                 },
                 "web_pool3": {
                     "class": "Pool",
                     "monitors": [
                         "http"
                     ],
                     "members": [
                         {
                             "servicePort": 80,
                             "serverAddresses": [
                                 "192.3.1.100",
                                 "192.3.1.110"
                             ]
                         }
                     ]
                 }
             }
         }
     }
 }

```

* `perApplication_example` - Per Application Example - JSON file with multiple Applications (and no Tenant Details)
 
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

## Import

As3 resources can be imported using the partition name, e.g., ( use comma separated partition names if there are multiple partitions in as3 deployments )

```
   terraform import bigip_as3.test Sample_http_01
   terraform import bigip_as3.test Sample_http_01,Sample_non_http_01
```

#### Import examples ( single and multiple partitions )

```

$ terraform import bigip_as3.test Sample_http_01
bigip_as3.test: Importing from ID "Sample_http_01"...
bigip_as3.test: Import prepared!
  Prepared bigip_as3 for import
bigip_as3.test: Refreshing state... [id=Sample_http_01]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

$ terraform show
# bigip_as3.test:
resource "bigip_as3" "test" {
    as3_json      = jsonencode(
        {
            action      = "deploy"
            class       = "AS3"
            declaration = {
                Sample_http_01 = {
                    A1    = {
                        class      = "Application"
                        jsessionid = {
                            class             = "Persist"
                            cookieMethod      = "hash"
                            cookieName        = "JSESSIONID"
                            persistenceMethod = "cookie"
                        }
                        service    = {
                            class              = "Service_HTTP"
                            persistenceMethods = [
                                {
                                    use = "jsessionid"
                                },
                            ]
                            pool               = "web_pool"
                            virtualAddresses   = [
                                "10.0.2.10",
                            ]
                        }
                        web_pool   = {
                            class    = "Pool"
                            members  = [
                                {
                                    serverAddresses = [
                                        "192.0.2.10",
                                        "192.0.2.11",
                                    ]
                                    servicePort     = 80
                                },
                            ]
                            monitors = [
                                "http",
                            ]
                        }
                    }
                    class = "Tenant"
                }
                class          = "ADC"
                id             = "UDP_DNS_Sample"
                label          = "UDP_DNS_Sample"
                remark         = "Sample of a UDP DNS Load Balancer Service"
                schemaVersion  = "3.0.0"
            }
            persist     = true
        }
    )
    id            = "Sample_http_01"
    tenant_filter = "Sample_http_01"
    tenant_list   = "Sample_http_01"
}





$ terraform import bigip_as3.test Sample_http_01,Sample_non_http_01
bigip_as3.test: Importing from ID "Sample_http_01,Sample_non_http_01"...
bigip_as3.test: Import prepared!
  Prepared bigip_as3 for import
bigip_as3.test: Refreshing state... [id=Sample_http_01,Sample_non_http_01]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

$ terraform show
# bigip_as3.test:
resource "bigip_as3" "test" {
    as3_json      = jsonencode(
        {
            action      = "deploy"
            class       = "AS3"
            declaration = {
                Sample_http_01     = {
                    A1    = {
                        class      = "Application"
                        jsessionid = {
                            class             = "Persist"
                            cookieMethod      = "hash"
                            cookieName        = "JSESSIONID"
                            persistenceMethod = "cookie"
                        }
                        service    = {
                            class              = "Service_HTTP"
                            persistenceMethods = [
                                {
                                    use = "jsessionid"
                                },
                            ]
                            pool               = "web_pool"
                            virtualAddresses   = [
                                "10.0.2.10",
                            ]
                        }
                        web_pool   = {
                            class    = "Pool"
                            members  = [
                                {
                                    serverAddresses = [
                                        "192.0.2.10",
                                        "192.0.2.11",
                                    ]
                                    servicePort     = 80
                                },
                            ]
                            monitors = [
                                "http",
                            ]
                        }
                    }
                    class = "Tenant"
                }
                Sample_non_http_01 = {
                    DNS_Service = {
                        Pool1   = {
                            class    = "Pool"
                            members  = [
                                {
                                    serverAddresses = [
                                        "10.1.10.100",
                                    ]
                                    servicePort     = 53
                                },
                                {
                                    serverAddresses = [
                                        "10.1.10.101",
                                    ]
                                    servicePort     = 53
                                },
                            ]
                            monitors = [
                                "icmp",
                            ]
                        }
                        class   = "Application"
                        service = {
                            class            = "Service_UDP"
                            pool             = "Pool1"
                            virtualAddresses = [
                                "10.1.20.121",
                            ]
                            virtualPort      = 53
                        }
                    }
                    class       = "Tenant"
                }
                class              = "ADC"
                id                 = "UDP_DNS_Sample"
                label              = "UDP_DNS_Sample"
                remark             = "Sample of a UDP DNS Load Balancer Service"
                schemaVersion      = "3.0.0"
            }
            persist     = true
        }
    )
    id            = "Sample_http_01,Sample_non_http_01"
    tenant_filter = "Sample_http_01,Sample_non_http_01"
    tenant_list   = "Sample_http_01,Sample_non_http_01"
}

```

* `AS3 documentation` - https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/composing-a-declaration.html