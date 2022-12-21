---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_server_ssl"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_server_ssl resource
---

# bigip_ltm_profile_server_ssl

`bigip_ltm_profile_server_ssl` Manages server SSL profiles on a BIG-IP

Resources should be named with their "full path". The full path is the combination of the partition + name (example: /Common/my-pool ) or  partition + directory + name of the resource  (example: /Common/test/my-pool )

## Example Usage
    

```hcl

resource "bigip_ltm_profile_server_ssl" "test-ServerSsl" {
  name          = "/Common/test-ServerSsl"
  defaults_from = "/Common/serverssl"
  authenticate  = "always"
  ciphers       = "DEFAULT"
}

```      

## Argument Reference

* `name` (Required,type `string`) Specifies the name of the profile.Name of Profile should be full path,full path is the combination of the `partition + profile name`. For example `/Common/test-serverssl-profile`.

* `defaults_from` - (Optional) The parent template of this monitor template. Once this value has been set, it cannot be changed. By default, this value is `/Common/serverssl`.

* `cert` - (Optional) Specifies the name of the certificate that the system uses for server-side SSL processing.

* `key` - (Optional) Specifies the file name of the SSL key.

* `chain` - (Optional) Specifies the certificates-key chain to associate with the SSL profile

* `ciphers` - (Optional) Specifies the list of ciphers that the system supports. When creating a new profile, the default cipher list is provided by the parent profile.

* `cipher_group` - (Optional) Specifies the cipher group for the SSL server profile. It is mutually exclusive with the argument, `ciphers`. The default value is `none`.

* `peer_cert_mode` - (Optional) Specifies the way the system handles client certificates.When ignore, specifies that the system ignores certificates from client systems.When require, specifies that the system requires a client to present a valid certificate.When request, specifies that the system requests a valid certificate from a client but always authenticate the client.

* `authenticate` - (Optional) Specifies the frequency of server authentication for an SSL session.When `once`,specifies that the system authenticates the server once for an SSL session.
When `always`, specifies that the system authenticates the server once for an SSL session and also upon reuse of that session.

* `tm_options` - (Optional,type `list`) List of Enabled selection from a set of industry standard options for handling SSL processing.By default,
Don't insert empty fragments and No TLSv1.3 are listed as Enabled Options. `Usage` : tm_options    = ["dont-insert-empty-fragments","no-tlsv1.3"]

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

* `ssl_c3d` (Optional) Enables or disables SSL forward proxy bypass on receiving
 handshake_failure, protocol_version or unsupported_extension alert message during the serverside SSL handshake. When enabled and there is an SSL handshake_failure, protocol_version or unsupported_extension alert during the serverside SSL handshake, SSL traffic bypasses the BIG-IP system untouched, without decryption/encryption. The default value is disabled. Conversely, you can specify enabled to use this feature.

* `c3d_ca_cert` (Optional) Specifies the name of the certificate file that is used as the certification authority certificate when SSL client certificate constrained delegation is enabled. The certificate should be generated and installed by you on the system. When selecting this option, type a certificate file name.

* `c3d_ca_key` (Optional) Specifies the name of the key file that is used as the certification authority key when SSL client certificate constrained delegation is enabled. The key should be generated and installed by you on the system. When selecting this option, type a key file name.

* `c3d-ca-passphrase` (Optional) Specifies the passphrase of the key file that is used as the certification authority key when SSL client certificate constrained delegation is enabled. When selecting this option, type the passphrase corresponding to the selected c3d-ca-key.

* `c3d-cert-extension-custom-oids` (Optional) Specifies the custom extension OID of the client certificates to be included in the generated certificates using SSL client certificate constrained delegation.

* `c3d_cert_extension_includes` (Optional) Specifies the extensions of the client certificates to be included in the generated certificates using SSL client certificate constrained delegation. For example, { basic-constraints }. The default value is { basic-constraints extended-key-usage key-usage subject-alternative-name }. The extensions are:

	    basic-constraints
		  Basic constraints are used to indicate whether the certificate belongs
      to a CA.

	    extended-key-usage
		  Extended Key Usage is used, typically on a leaf certificate, to 
      indicate the purpose of the public key contained in the certificate.

	    key-usage
		  Key Usage provides a bitmap specifying the cryptographic operations 
      which may be performed using the public key contained in the 
      certificate; for example, it could indicate that the key should be 
      used for signature but not for encipherment.

	    subject-alternative-name
		  Subject Alternative Name allows identities to be bound to the subject
      of the certificate. These identities may be included in addition to 
      or in place of the identity in the subject field of the certificate.

* `c3d-cert-lifespan` Specifies the lifespan of the certificate generated using the SSL client certificate constrained delegation. The default value is 24.

## Importing
An existing server-ssl profile can be imported into this resource by supplying server-ssl profile Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_profile_server_ssl.test-ServerSsl-import /Common/test-ServerSsl

```
