---
layout: "bigip"
page_title: "Provider: BIG-IP"
description: |-
  Provides details about BIG-IP Terraform Provider
---

# F5 BIG-IP Terraform Provider

Use the F5 BIG-IP Terraform Provider to manage and provision your BIG-IP
configurations in Terraform. Using BIG-IP Provider you can manage LTM(Local Traffic Manager),Network,System objects and it also supports AS3/DO integration.

### Requirements

This provider uses the iControlREST API. All the resources are validated with BigIP v12.1.1 and above.

~> **NOTE** For AWAF resources, F5 BIG-IP version should be > v16.x , and ASM need to be provisioned.

## Example Usage
```hcl
variable hostname {}
variable username {}
variable password {}

terraform {
  required_providers {
    bigip = {
      source = "terraform-providers/bigip"
    }
  }
  required_version = ">= 0.13"
}

provider "bigip" {
  address  = var.hostname
  username = var.username
  password = var.password
}
```

## Argument Reference

- `address` - (type `string`) Domain name or IP address of the BIG-IP. Can be set via the `BIGIP_HOST` environment variable.
- `username` - (type `string`) BIG-IP Username for authentication. Can be set via the `BIGIP_USER` environment variable.
- `password` - (type `string`) BIG-IP Password for authentication. Can be set via the `BIGIP_PASSWORD` environment variable.
- `token_auth` - (Optional, Default `false`) Enable to use an external authentication source (LDAP, TACACS, etc). Can be set via the `BIGIP_TOKEN_AUTH` environment variable.
- `token_value` - (Optional) A token generated outside the provider, in place of password
- `login_ref` - (Optional,Default `tmos`) Login reference for token authentication (see BIG-IP REST docs for details). May be set via the `BIGIP_LOGIN_REF` environment variable.
- `port` - (Optional) Management Port to connect to BIG-IP,this is mainly required if we have single nic BIG-IP in AWS/Azure/GCP (or) Management port other than `443`. Can be set via `BIGIP_PORT` environment variable.
- `validate_certs_disable` - (Optional, Default `true`) If set to true, Disables TLS certificate check on BIG-IP. Can be set via the `BIGIP_VERIFY_CERT_DISABLE` environment variable.
- `trusted_cert_path` - (type `string`) Provides Certificate Path to be used TLS Validate.It will be required only if `validate_certs_disable` set to `false`.Can be set via the `BIGIP_TRUSTED_CERT_PATH` environment variable.

~> **Note** For BIG-IQ resources these provider credentials `address`,`username`,`password` can be set to BIG-IQ credentials.

~> **Note** The F5 BIG-IP provider gathers non-identifiable usage data for the purposes of improving the product as outlined in the end user license agreement for BIG-IP. To opt out of data collection, use the following : `export TEEM_DISABLE=true`
