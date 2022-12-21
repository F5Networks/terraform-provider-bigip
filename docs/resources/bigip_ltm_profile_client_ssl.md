---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_client_ssl"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_client_ssl resource
---

# bigip_ltm_profile_client_ssl

`bigip_ltm_profile_client_ssl` Manages client SSL profiles on a BIG-IP

Resources should be named with their "full path". The full path is the combination of the partition + name (example: /Common/my-pool ) or  partition + directory + name of the resource  (example: /Common/test/my-pool )

## Example Usage
    
```hcl
resource "bigip_ltm_profile_client_ssl" "test-ClientSsl" {
  name          = "/Common/test-ClientSsl"
  defaults_from = "/Common/clientssl"
  authenticate  = "always"
  ciphers       = "DEFAULT"
}
```      

## Argument Reference

* `name` (Required,type `string`) Specifies the name of the profile.Name of Profile should be full path.The full path is the combination of the `partition + profile name`,For example `/Common/test-clientssl-profile`.

* `defaults_from` - (Optional) Parent profile for this clientssl profile.Once this value has been set, it cannot be changed. Default value is `/Common/clientssl`. It Should Full path `/partition/profile_name`

* `allow_non_ssl` - (Optional) Enables or disables acceptance of non-SSL connections, When creating a new profile, the setting is provided by the parent profile

* `authenticate` - (Optional) Specifies the frequency of client authentication for an SSL session.When `once`,specifies that the system authenticates the client once for an SSL session.
When `always`, specifies that the system authenticates the client once for an SSL session and also upon reuse of that session.

* `tm_options` - (Optional,type `list`) List of Enabled selection from a set of industry standard options for handling SSL processing.By default,
Don't insert empty fragments and No TLSv1.3 are listed as Enabled Options. `Usage` : tm_options    = ["dont-insert-empty-fragments","no-tlsv1.3"]

* `authenticate_depth` - (Optional) Specifies the maximum number of certificates to be traversed in a client certificate chain

* `cert` - (Optional) Specifies a cert name for use.

* `key` - (Optional) Contains a key name

* `chain` - (Optional) Contains a certificate chain that is relevant to the certificate and key mentioned earlier.This key is optional

* `ciphers` - (Optional) Specifies the list of ciphers that the system supports. When creating a new profile, the default cipher list is provided by the parent profile.

* `cipher_group` - (Optional) Specifies the cipher group for the SSL server profile. It is mutually exclusive with the argument, `ciphers`. The default value is `none`.

* `peer_cert_mode` - (Optional) Specifies the way the system handles client certificates.When ignore, specifies that the system ignores certificates from client systems.When require, specifies that the system requires a client to present a valid certificate.When request, specifies that the system requests a valid certificate from a client but always authenticate the client.

* `renegotiation` - (Optional) Enables or disables SSL renegotiation.When creating a new profile, the setting is provided by the parent profile

* `retain_certificate` - (Optional) When `true`, client certificate is retained in SSL session.

* `secure_renegotiation` - (Optional) Specifies the method of secure renegotiations for SSL connections. When creating a new profile, the setting is provided by the parent profile.
When `request` is set the system request secure renegotation of SSL connections.
`require` is a default setting and when set the system permits initial SSL handshakes from clients but terminates renegotiations from unpatched clients.
The `require-strict` setting the system requires strict renegotiation of SSL connections. In this mode the system refuses connections to insecure servers, and terminates existing SSL connections to insecure servers

* `server_name` - (Optional) Specifies the fully qualified DNS hostname of the server used in Server Name Indication communications. When creating a new profile, the setting is provided by the parent profile.The server name can also be a wildcard string containing the asterisk `*` character.

* `sni_default` - (Optional) Indicates that the system uses this profile as the default SSL profile when there is no match to the server name, or when the client provides no SNI extension support.When creating a new profile, the setting is provided by the parent profile.
There can be only one SSL profile with this setting enabled.

* `sni_require` - (Optional) Requires that the network peers also provide SNI support, this setting only takes effect when `sni_default` is set to `true`.When creating a new profile, the setting is provided by the parent profile

* `strict_resume` - (Optional) Enables or disables the resumption of SSL sessions after an unclean shutdown.When creating a new profile, the setting is provided by the parent profile. 

* `ssl_forward_proxy` - (Optional) Specifies whether SSL forward proxy feature is enabled or not. The default value is disabled.

* `ssl_forward_proxy_bypass` - (Optional) Specifies whether SSL forward proxy bypass feature is enabled or not. The default value is disabled.

* `ssl_c3d` (Optional) Enables or disables SSL client certificate constrained delegation. The default option is disabled. Conversely, you can specify enabled to use the SSL client certificate constrained delegation.
  
* `c3d_client_fallback_cert` (Optional) Specifies the client certificate to use in SSL client certificate constrained delegation. This certificate will be used if client does not provide a cert during the SSL handshake. The default value is none.

* `c3d_drop_unknown_ocsp_status` (Optional) Specifies the BIG-IP action when the OCSP responder returns unknown status. The default value is drop, which causes the onnection to be dropped. Conversely, you can specify ignore, which causes the connection to ignore the unknown status and continue.

* `c3d_ocsp` (Optional) Specifies the SSL client certificate constrained delegation OCSP object that the BIG-IP SSL should use to connect to the OCSP responder and check the client certificate status.


## Importing
An existing client-ssl profile can be imported into this resource by supplying client-ssl profile Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_profile_client_ssl.test-ClientSsl-import /Common/test-ClientSsl
```
