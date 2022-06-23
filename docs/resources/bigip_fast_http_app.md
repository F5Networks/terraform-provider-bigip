---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_http_app"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_fast_http_app resource
---

# bigip_fast_tcp_app

`bigip_fast_http_app` This resource will create and manage FAST HTTP applications on BIG-IP from provided JSON declaration. 


## Example Usage


```hcl

resource "bigip_fast_http_app" "fast-http" {
  tenant = "httptenant"
  application= "httpapptest"
  virtual_server = {
    ip = "10.200.20.1"
    port = 201
  }
}

```      

## Argument Reference

* `application` - (Required) Name of the FAST HTTP application.

* `tenant` - (Required) Name of the FAST HTTP application tenant.

* `virtual_server.ip` - (Optional) This IP address, combined with the port you specify below, becomes the BIG-IP virtual server address and port, which clients use to access the application.

* `virtual_server.port` - (Optional) This is the virtual server port address, which is combined with the virtual IP address to access the application. 


## Attributes Reference

* `fast_tcp_json` - Exported json containing properties of the deployed FAST TCP application.****

* `FAST documentation` - https://clouddocs.f5.com/products/extensions/f5-appsvcs-templates/latest/
