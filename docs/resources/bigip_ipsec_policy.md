---
layout: "bigip"
page_title: "BIG-IP: bigip_ipsec_policy"
sidebar_current: "docs-bigip-resource-ipsec-policy-x"
description: |-
   Provides details about bigip_ipsec_policy resource
---

# bigip_ipsec_policy

`bigip_ipsec_policy` Manage IPSec policies on a BIG-IP

Resources should be named with their "full path". The full path is the combination of the partition + name (example: /Common/test-policy)


## Example Usage

```hcl
resource "bigip_ipsec_policy" "test-policy" {
  name                  = "/Common/test-policy"
  description           = "created by terraform provider"
  protocol              = "esp"
  mode                  = "tunnel"
  tunnel_local_address  = "192.168.1.1"
  tunnel_remote_address = "10.10.1.1"
  auth_algorithm        = "sha1"
  encrypt_algorithm     = "3des"
  lifetime              = 3
  ipcomp                = "deflate"
}
```      

## Argument Reference
* `name` - (Required) Name of the IPSec policy,it should be "full path".The full path is the combination of the partition + name of the IPSec policy.(For example `/Common/test-policy`)

* `description` - (Optional,type `string`) Description of the IPSec policy.

* `protocol` - (Optional,type `string`) Specifies the IPsec protocol. Valid choices are: `ah, esp` 

* `mode` - (Optional,type `string`) Specifies the processing mode. Valid choices are: `transport, interface, isession, tunnel`

* `tunnel_local_address` - (Optional,type `string`) Specifies the local endpoint IP address of the IPsec tunnel. This parameter is only valid when mode is tunnel.

* `tunnel_remote_address` - (Optional, type `string`) Specifies the remote endpoint IP address of the IPsec tunnel. This parameter is only valid when mode is tunnel.

* `encrypt_algorithm` - (Optional, type `string`) Specifies the algorithm to use for IKE encryption. Valid choices are: `null, 3des, aes128, aes192, aes256, aes-gmac256,
  aes-gmac192, aes-gmac128, aes-gcm256, aes-gcm192, aes-gcm256, aes-gcm128`

* `auth_algorithm` - (Optional, type `string`) Specifies the algorithm to use for IKE authentication. Valid choices are: `sha1, sha256, sha384, sha512, aes-gcm128,
  aes-gcm192, aes-gcm256, aes-gmac128, aes-gmac192, aes-gmac256`

* `lifetime` - (Optional, type `int`) Specifies the length of time before the IKE security association expires, in minutes.

* `kb_lifetime` - (Optional, type `int`) Specifies the length of time before the IKE security association expires, in kilobytes.

* `perfect_forward_secrecy` - (Optional, type `string`) Specifies the Diffie-Hellman group to use for IKE Phase 2 negotiation. Valid choices are: `none, modp768, modp1024, modp1536, modp2048, modp3072,
  modp4096, modp6144, modp8192`

* `ipcomp` - (Optional, type `string`) Specifies whether to use IPComp encapsulation. Valid choices are: `none", null", deflate`
