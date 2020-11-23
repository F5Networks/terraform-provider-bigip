---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_ssl_certificate"
sidebar_current: "docs-bigip-datasource-ssl_certificate-x"
description: |-
    Provides details about bigip_ltm_ssl_certificate data source
---

# bigip\_ltm\_ssl_certificate

Use this data source (`bigip_ltm_ssl_certificate`) to get the ltm irule details available on BIG-IP
 
 
## Example Usage
```hcl

data "bigip_ltm_ssl_certificate" "test" {
  name = "terraform_ssl_certificate"
  partition = "Common"
}


output "bigip_ssl_certificate" {
  value = "${data.bigip_ltm_ssl_certificate.test.irule}"
}

```      

## Argument Reference

* `name` - (Required) Name of the ssl_certificate

* `partition` - (Required) partition of the ltm ssl_certificate


## Attributes Reference

* `ssl_certificate` - Irule configured on bigip

* `name` - Name of ssl_certificate configured on bigip with full path

* `partition` - Bigip partition in which rule is configured

