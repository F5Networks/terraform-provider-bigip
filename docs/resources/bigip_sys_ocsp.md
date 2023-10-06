---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_ocsp"
subcategory: "System"
description: |-
  Provides details about OCSP resource for BIG-IP
---

# bibip\_sys\_ocsp

`bigip_sys_ocsp` Manages F5 BIG-IP OCSP responder using iControl REST.

## Example Usage

```hcl
resource "bigip_sys_ocsp" "test-ocsp" {
  name              = "/Uncommon/test-ocsp"
  proxy_server_pool = "/Common/test-poolxyz"
  signer_key        = "/Common/le-ssl"
  signer_cert       = "/Common/le-ssl"
  passphrase        = "testabcdef"
}
```

## Argument Reference

* `name` - (Required,type `string`) Name of the OCSP Responder. Name should be in pattern `/partition/ocsp_name`.

* `proxy_server_pool` - (Required,type `string`) Specifies the proxy server pool the BIG-IP system uses to fetch the OCSP response.

* `dns_resolver` - (Optional,type `string`) Specifies the internal DNS resolver the BIG-IP system uses to fetch the OCSP response.

* `route_domain` - (Optional,type `string`) Specifies the route domain for the OCSP responder.

* `concurrent_connections_limit` - (Optional,type `int`) Specifies the maximum number of connections per second allowed for the OCSP certificate validator. The default value is `50`.

* `responder_url` - (Optional,type `string`) Specifies the URL of the OCSP responder.

* `connection_timeout` - (Optional,type `int`) Specifies the time interval that the BIG-IP system waits for before ending the connection to the OCSP responder, in seconds. The default value is `8`.

* `trusted_responders` - (Optional,type `string`) Specifies the certificates used for validating the OCSP response.

* `clock_skew` - (Optional,type `int`) Specifies the time interval that the BIG-IP system allows for clock skew, in seconds. The default value is `300`.

* `status_age` - (Optional,type `int`) Specifies the maximum allowed lag time that the BIG-IP system accepts for the 'thisUpdate' time in the OCSP response, in seconds. The default value is `0`.

* `strict_resp_cert_check` - (Optional,type `string`) Specifies whether the responder's certificate is checked for an OCSP signing extension. The default value is `enabled`.

* `cache_timeout` - (Optional,type `string`) Specifies the lifetime of the OCSP response in the cache, in seconds. The default value is `indefinite`.

* `cache_error_timeout` - (Optional,type `string`) Specifies the lifetime of an error response in the cache, in seconds. This value must be greater than connection_timeout. The default value is `3600`.

* `signer_cert` - (Required,type `string`) Specifies the certificate used to sign the OCSP request.

* `signer_key` - (Required,type `string`) Specifies the key used to sign the OCSP request.

* `passphrase` - (Optional,type `string`) Specifies a passphrase used to sign an OCSP request.

* `sign_hash` - (Optional,type `string`) Specifies the hash algorithm used to sign the OCSP request. The default value is `sha256`.


## Importing
An existing OCSP can be imported into this resource by supplying the full path name  ex : `/partition/name`
An example is below:
```sh
$ terraform import bigip_sys_ocsp.test-ocsp /Common/test-ocsp
```
