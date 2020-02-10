---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_application"
sidebar_current: "docs-bigip-datasource-application-x"
description: |-
   Provides details about bigip_as3_application datasource
---

# bigip\_as3\_application

`bigip_as3_application` Manages an Application class, which comprises a set of resources used to manage, secure, and enhance the delivery of a simple or complex network-based application

The basic resources are virtual servers, profiles, iRules, pools, pool members, and monitors. At a minimum, you must include the application type. In the following example, App1 is the name of the application.

## Example Usage


```hcl
data "bigip_as3_app" "App1" {
  name = "App1"
  template = "https"
  pool_class = "${data.bigip_as3_pool.mydataas3pool.id}"
  service_class = "${data.bigip_as3_service.myservice.id}"
  cert_class = "${data.bigip_as3_cert.exmpcert.id}"
  tls_server_class = "${data.bigip_as3_tls_server.exmpserver.id}"
  enable = true
}
```

## Argument Reference

* `name` - (Required) Name of the application

* `template` - (Required) Each application type specified has certain required and default elements and selects appropriate setup of various ADC/Security features

* `enable` - (Optional) Application handles traffic only when enabled(default) 

* `label` - (Optional) Optional friendly name for this object

* `remark` - (Optional) Arbitrary (brief) text pertaining to this object

* `schema_overlay` - (Optional) BIG-IQ name for a supplemental validation schema is applied to the Application class definition before the main AS3 schema

* `servicemain`- (Optional) Primary service of the application

* `pool_class` - (Optional) Pointer to the id of pool datasource

* `service_class` - (Optional) Pointer to the id of service datasource

* `cert_class` - (Optional) Pointer to the id of certificate datasource

* `tls_server_class` - (Optional) Pointer to the id of TLS Server datasource
