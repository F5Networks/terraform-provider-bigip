---
layout: "bigip"
page_title: "BIG-IP: bigip_ssl_certificate"
subcategory: "System"
description: |-
  Provides details about bigip_ssl_certificate resource
---

# bigip_ssl_certificate

`bigip_ssl_certificate` This resource will import SSL certificates on BIG-IP LTM. 
Certificates can be imported from certificate files on the local disk, in PEM format


## Example Usage


```hcl

resource "bigip_ssl_certificate" "test-cert" {
  name      = "servercert.crt"
  content   = file("servercert.crt")
  partition = "Common"
}

```      

## Argument Reference


* `name`- (Required) Name of the SSL Certificate to be Imported on to BIGIP

* `content` - (Required) Content of certificate on Local Disk,path of SSL certificate will be provided to terraform `file` function 

* `partition` - Partition on to SSL Certificate to be imported. The parameter is not required when running terraform import operation. In such case the name must be provided in full_path format.

* `monitoring_type` - Specifies the type of monitoring used.

* `issuer_cert` - Specifies the issuer certificate.

* `ocsp` - Specifies the OCSP responder.
