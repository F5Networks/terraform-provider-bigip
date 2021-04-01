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


# Example Usage for json file
resource "bigip_as3"  "as3-example1" {
       as3_json = file("example1.json")
 }

# Example Usage for json file with tenant filter
resource "bigip_as3"  "as3-example1" {
       as3_json = file("example2.json")
       tenant_filter = "Sample_03"
 }


```

## Argument Reference


* `as3_json` - (Required) Path/Filename of Declarative AS3 JSON which is a json file used with builtin ```file``` function

* `tenant_filter` - (Optional) If there are multiple tenants on a BIG-IP, this attribute helps the user to set a particular tenant to which he want to reflect the changes. Other tenants will neither be created nor be modified.

* `ignore_metadata` - (Optional) Set True if you want to ignore metadata changes during update. By default it is set to false

* `as3_example1.json` - Example  AS3 Declarative JSON file with single tenant

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

