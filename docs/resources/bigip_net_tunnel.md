---
layout: "bigip"
page_title: "BIG-IP: bigip_net_tunnel"
subcategory: "Network"
description: |-
  Provides details about bigip_net_tunnel resource
---

# bigip\_net\_tunnel

`bigip_net_tunnel` Manages a tunnel configuration



## Example Usage


```hcl
resource "bigip_net_tunnel" "example1" {
  name          = "example1"
  local_address = "192.16.81.240"
  profile       = "/Common/dslite"
}

```

## Argument Reference

* `name` - (Required) Name of the tunnel

* `local_address` - (Required) Specifies a local IP address. This option is required

* `profile` - (Required) Specifies the profile that you want to associate with the tunnel     

* `app_service` - (Optional) The application service that the object belongs to

* `auto_last_hop` - (Optional) Specifies whether auto lasthop is enabled or not

* `description` - (Optional) User defined description

* `mode` - (Optional) Specifies how the tunnel carries traffic

* `partition` - (Optional) Displays the admin-partition within which this component resides

* `remote_address` - (Optional) Specifies a remote IP address

* `secondary_address` - (Optional) Specifies a secondary non-floating IP address when the local-address is set to a floating address

* `tos` - (Optional) Specifies a value for insertion into the Type of Service (ToS) octet within the IP header of the encapsulating header of transmitted packets

* `traffic_group` - (Optional) Specifies a traffic-group for use with the tunnel

* `transparent` - (Optional) Enables or disables the tunnel to be transparent

* `use_pmtu` - (Optional) Enables or disables the tunnel to use the PMTU (Path MTU) information provided by ICMP NeedFrag error messages

* `idle_timeout` - (Optional) Specifies an idle timeout for wildcard tunnels in seconds

* `key` - (Optional) The key field may represent different values depending on the type of the tunnel

* `mtu` - (Optional) Specifies the maximum transmission unit (MTU) of the tunnel
