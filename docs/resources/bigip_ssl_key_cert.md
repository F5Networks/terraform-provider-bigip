---
layout: "bigip"
page_title: "BIG-IP: bigip_ssl_key_cert"
subcategory: "System"
description: |-
  Provides details about bigip_ssl_key_cert resource
---

# bigip_ssl_key_cert

`bigip_ssl_key_cert` This resource will import SSL certificate and key on BIG-IP LTM. 
The certificate and the key can be imported from files on the local disk, in PEM format


## Example Usage


```hcl

resource "bigip_ssl_key_cert" "testkeycert" {
  partition    = "Common"
  key_name     = "ssl-test-key"
  key_content  = file("key.pem")
  cert_name    = "ssl-test-cert"
  cert_content = file("certificate.pem")
}

```      

## Argument Reference


* `key_name`- (Required,type `string`) Name of the SSL key to be Imported on to BIGIP.

* `key_content` - (Required) Content of SSL key on Local Disk,path of SSL key will be provided to terraform `file` function.

* `cert_name`- (Required,type `string`) Name of the SSL certificate to be Imported on to BIGIP.

* `cert_content` - (Required) Content of certificate on Local Disk,path of SSL certificate will be provided to terraform `file` function.

* `partition` - (Optional,type `string`) Partition on to SSL certificate and key to be imported.

* `passphrase` - (Optional,type `string`) Passphrase on the SSL key.

* `cert_monitoring_type` - (Optional,type `string`) Specifies the type of monitoring used.

* `issuer_cert` - (Optional,type `string`) Specifies the issuer certificate.

* `cert_ocsp` - (Optional,type `string`) Specifies the OCSP responder.


## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - identifier of the resource.

* `key_full_path` - full path of the SSL key on the BIGIP.

* `cert_full_path` - full path of the SSL certificate on the BIGIP.
