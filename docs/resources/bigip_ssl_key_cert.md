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

* `key_content_wo` - (Optional) Write-only alternative to `key_content`. Terraform sends this key content to the provider but does not persist it in plan or state. Cannot be used with `key_content`.

* `key_content_wo_version` - (Optional) A write-only version marker for `key_content_wo`. Change this value to trigger a key re-upload.

* `cert_name`- (Required,type `string`) Name of the SSL certificate to be Imported on to BIGIP.

* `cert_content` - (Required) Content of certificate on Local Disk,path of SSL certificate will be provided to terraform `file` function.

* `cert_content_wo` - (Optional) Write-only alternative to `cert_content`. Terraform sends this certificate content to the provider but does not persist it in plan or state. Cannot be used with `cert_content`.

* `cert_content_wo_version` - (Optional) A write-only version marker for `cert_content_wo`. Change this value to trigger a certificate re-upload.

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
