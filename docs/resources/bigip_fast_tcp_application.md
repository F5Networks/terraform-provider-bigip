---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_tcp_application"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_fast_tcp_application resource
---

# bigip_fast_tcp_application

`bigip_fast_tcp_application` This resource will create and manage FAST TCP applications on BIG-IP from provided JSON declaration. 


## Example Usage


```hcl

resource "bigip_fast_tcp_application" "fast-tcp-app" {
  application = "tcp_app_2"
  tenant      = "tcp_app_tenant"
  virtual_server = {
    ip = "11.12.16.30"
    port = 443
  }
  fastl4 = {
    enable = true
    generate_fastl4_profile = false
    fastl4_profile_name = "/Common/apm-forwarding-fastL4"
  }
}

```      

## Argument Reference

* `application` - (Required) Name of the FAST TCP application.
* `tenant` - (Required) Name of the FAST TCP application tenant.
* `virtual_server.ip` - (Optional) This IP address, combined with the port you specify below, becomes the BIG-IP virtual server address and port, which clients use to access the application.
* `virtual_server.port` - (Optional) This is the virtual server port address, which is combined with the virtual IP address to access the application. 
* `fastl4.enable` - (Optional) Determines whether to use fastl4 protocol profiles or not.
* `fastl4.generate_fastl4_profile` - (Optional) Determines whether to use FAST generated fastl4 protocol profiles.
* `fastl4.fastl4_profile_name` - (Optional) Name of an existing fastl4 protocol profile.


## Attributes Reference

* `fast_tcp_json` - Exported json containing properties of the deployed FAST TCP application.****
