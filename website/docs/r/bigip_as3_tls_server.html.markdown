---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_tls_server"
sidebar_current: "docs-bigip-datasource-tls-server-x"
description: |-
   Provides details about bigip_as3_tls_server datasource
---

# bigip\_as3\_tls\_server

`bigip_as3_tls_server` Manages a TLS Server class, which contain server parameters (connections arriving to ADC)

## Example Usage


```hcl
data "bigip_as3_tls_server" "exmpserver" {
  name = "exmpserver"
  certificates {
    certificate = "exmpcert"
  }
}
```

## Argument Reference

* `name` - (Required) Name of the tls server

* Below attributes needs to be configured under certificates option.

* `certificate` - (Optional) Specifies the name of the certificate datasource
