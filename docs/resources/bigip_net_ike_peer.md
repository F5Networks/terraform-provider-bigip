---
layout: "bigip"
page_title: "BIG-IP: bigip_net_ike_peer"
subcategory: "Network"
description: |-
  Provides details about bigip_net_ike_peer resource
---

# bigip\_net\_ike_peer

`bigip_net_ike_peer` Manages a ike_peer configuration



## Example Usage


```hcl
resource "bigip_net_ike_peer" "example1" {
  name          = "example1"
  local_address = "192.16.81.240"
  profile       = "/Common/dslite"
}
```

## Argument Reference

* `name` - (Required) Name of the ike_peer

* `app_service` - (Optional)The application service that the object belongs to 

* `ca_cert_file` - (Optional)the trusted root and intermediate certificate authorities 

* `crl_file` - (Optional)Specifies the file name of the Certificate Revocation List. Only supported in IKEv1 

* `description` - (Optional)User defined description 

* `generate_policy` - (Optional)Enable or disable the generation of Security Policy Database entries(SPD) when the device is the responder of the IKE remote node 

* `mode` - (Optional)Defines the exchange mode for phase 1 when racoon is the initiator, or the acceptable exchange mode when racoon is the responder

* `my_cert_file` - (Optional)Specifies the name of the certificate file object

* `my_cert_key_file` - (Optional)Specifies the name of the certificate key file object 

* `my_cert_key_passphrase` - (Optional)Specifies the passphrase of the key used for my-cert-key-file 

* `my_id_type` - (Optional)Specifies the identifier type sent to the remote host to use in the phase 1 negotiation 

* `my_id_value` - (Optional)Specifies the identifier value sent to the remote host in the phase 1 negotiation 

* `nat_traversal` - (Optional)Enables use of the NAT-Traversal IPsec extension 

* `passive` - (Optional)Specifies whether the local IKE agent can be the initiator of the IKE negotiation with this ike-peer

* `peers_cert_file` - (Optional)Specifies the peer’s certificate for authentication 

* `peers_cert_type` - (Optional)Specifies that the only peers-cert-type supported is certfile

* `peers_id_type` - (Optional)Specifies which of address, fqdn, asn1dn, user-fqdn or keyid-tag types to use as peers-id-type 

* `peers_id_value` - (Optional)Specifies the peer’s identifier to be received 

* `phase1_auth_method` - (Optional)Specifies the authentication method used for phase 1 negotiation 

* `phase1_encrypt_algorithm` - (Optional)Specifies the encryption algorithm used for the isakmp phase 1 negotiation 

* `phase1_hash_algorithm` - (Optional)Defines the hash algorithm used for the isakmp phase 1 negotiation 

* `phase1_perfect_forward_secrecy` - (Optional)Defines the Diffie-Hellman group for key exchange to provide perfect forward secrecy 

* `preshared_key` - (Optional)Specifies the preshared key for ISAKMP SAs

* `preshared_key_encrypted` - (Optional)Display the encrypted preshared-key for the IKE remote node 

* `prf` - (Optional) Specifies the pseudo-random function used to derive keying material for all cryptographic operations

* `proxy_support` - (Optional)If this value is enabled, both values of ID payloads in the phase 2 exchange are used as the addresses of end-point of IPsec-SAs 

* `remote_address` - (Required)Specifies the IP address of the IKE remote node 

* `state` - (Optional)Enables or disables this IKE remote node 

* `traffic_selector` - (Optional)Specifies the names of the traffic-selector objects associated with this ike-peer 

* `verify_cert` - (Optional)Specifies whether to verify the certificate chain of the remote peer based on the trusted certificates in ca-cert-file 

* `version` - (Optional)Specifies which version of IKE to be used 

* `dpd_delay` - (Optional)Specifies the number of seconds between Dead Peer Detection messages 

* `lifetime` - (Optional)Defines the lifetime in minutes of an IKE SA which will be proposed in the phase 1 negotiations 

* `replay_window_size` - (Optional)Specifies the replay window size of the IPsec SAs negotiated with the IKE remote node 
