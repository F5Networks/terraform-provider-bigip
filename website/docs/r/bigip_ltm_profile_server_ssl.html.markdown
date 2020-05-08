---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_server_ssl"
sidebar_current: "docs-bigip-resource-profile_server_ssl-x"
description: |-
    Provides details about bigip_ltm_profile_server_ssl resource
---

# bigip_ltm_profile_server_ssl

`bigip_ltm_profile_server_ssl` Manages server SSL profiles on a BIG-IP



## Example Usage
    

```hcl

resource "bigip_ltm_profile_server_ssl" "test-ServerSsl" {
  name          = "/Common/test-ServerSsl"
  partition     = "Common"
  defaults_from = "/Common/serverssl"
  authenticate  = "always"
  ciphers       = "DEFAULT"
}
  
```      

## Argument Reference

* `name` (Required) Specifies the name of the profile. (type `string`)

* `partition` - (Optional) Device partition to manage resources on.

* `defaults_from` - (Optional) The parent template of this monitor template. Once this value has been set, it cannot be changed. By default, this value is `/Common/serverssl`.

* `cert` - (Optional) Specifies the name of the certificate that the system uses for server-side SSL processing.

* `key` - (Optional) Specifies the file name of the SSL key.

* `chain` - (Optional) Specifies the certificates-key chain to associate with the SSL profile

* `ciphers` - (Optional) Specifies the list of ciphers that the system supports. When creating a new profile, the default cipher list is provided by the parent profile.

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
