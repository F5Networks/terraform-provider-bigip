---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_service"
sidebar_current: "docs-bigip-datasource-service-x"
description: |-
   Provides details about bigip_as3_service datasource
---

# bigip\_as3\_service

`bigip_as3_service` Manages a Service class, which specifies each service and associated virtual IP address (called a virtual server on the BIG-IP system). Clients use the virtual IP address to access resources behind the BIG-IP system

If the template you specified in the Application class is http, https, tcp, udp, or l4, you MUST specify an object with the matching service class Service_HTTP, Service_HTTPS, Service_TCP, Service_UDP, or Service_L4 and name it serviceMain

## Example Usage


```hcl
data "bigip_as3_service" "myservice" {
  name = "serviceMain"
  virtual_addresses=["10.0.10.10"]
  pool_name = "${data.bigip_as3_pool.mydataas3pool.name}"
  server_tls = "${data.bigip_as3_tls_server.exmpserver.name}"
  service_type = "https"
  persistence_methods = ["cookie"]
}
```

## Argument Reference

* `name` - (Required) Name of the service

* `virtual_addresses` - (Optional) Virtual server will listen to each IP address in list.

* `pool_name` - (Required) Specifies the name of the pool related to the virtual server

* `service_type` - (Optional) Specifies the servive type that is being created 

* `persistence_methods` - (Optional) Default ‘cookie’ is generally good. Use ‘persistenceMethods: []’ for no persistence

* Below attribute is set only if the service type is https

* `server_tls` - (Optional) Specifies the name of the tls server related to the virtual server
