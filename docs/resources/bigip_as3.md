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
       as3_json = "${file("example1.json")}"
 }

# Example Usage for json file with tenant filter
resource "bigip_as3"  "as3-example1" {
       as3_json = "${file("example2.json")}"
       tenant_filter = "Sample_03"
 }


```

## Argument Reference


* `as3_json` - (Required) Path/Filename of Declarative AS3 JSON which is a json file used with builtin ```file``` function

* `tenant_filter` - (Optional) If there are muntiple tenants in a json this attribute helps the user to set a particular tenant to which he want to reflect the changes. Other tenants will neither be created nor be modified 

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

* `AS3 documentation` - https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/composing-a-declaration.html

