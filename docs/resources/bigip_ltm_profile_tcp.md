---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_tcp"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_tcp resource
---

# bigip\_ltm\_profile_tcp

`bigip_ltm_profile_tcp` Configures a custom TCP LTM Profile for use by health checks.

Resources should be named with their `full path`. The full path is the combination of the `partition + name` (example: /Common/my-pool ) or  `partition + directory + name` of the resource  (example: /Common/test/my-pool )

## Example Usage

```hcl
resource "bigip_ltm_profile_tcp" "sanjose-tcp-lan-profile" {
  name               = "/Common/sanjose-tcp-lan-profile"
  idle_timeout       = 200
  close_wait_timeout = 5
  finwait_2timeout   = 5
  finwait_timeout    = 300
  keepalive_interval = 1700
  deferred_accept    = "enabled"
  fast_open          = "enabled"
}
```      

## Argument Reference
	
* `name` (Required,type `string`) Name of the LTM TCP Profile,name should be `full path`. The full path is the combination of the `partition + name` (example: /Common/my-pool ) or  `partition + directory + name` of the resource  (example: /Common/test/my-pool )

* `defaults_from` - (Optional,type `string`) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `idle_timeout` - (Optional,type `int`) Specifies the number of seconds that a connection is idle before the connection is eligible for deletion. The default value is 300 seconds.

* `close_wait_timeout` - (Optional,type `int`) Specifies the number of seconds that a connection remains in a LAST-ACK state before quitting. A value of 0 represents a term of forever (or until the maxrtx of the FIN state). The default value is 5 seconds.

* `finwait_timeout` - (Optional,type `int`) Specifies the number of seconds that a connection is in the FIN-WAIT-1 or closing state before quitting. The default value is 5 seconds. A value of 0 (zero) represents a term of forever (or until the maxrtx of the FIN state). You can also specify immediate or indefinite.

* `finwait_2timeout` - (Optional,type `int`) Specifies the number of seconds that a connection is in the FIN-WAIT-2 state before quitting. The default value is 300 seconds. A value of 0 (zero) represents a term of forever (or until the maxrtx of the FIN state).

* `keepalive_interval` - (Optional,type `int`) Specifies the keep alive probe interval, in seconds. The default value is 1800 seconds.

* `zerowindow_timeout` - (Optional,type `int`) Specifies the timeout in milliseconds for terminating a connection with an effective zero length TCP transmit window.

* `send_buffersize` - (Optional,type `int`) Specifies the SEND window size. The default is 131072 bytes.

* `receive_windowsize` - (Optional,type `int`) Specifies the maximum advertised RECEIVE window size. This value represents the maximum number of bytes to which the RECEIVE window can scale. The default is 65535 bytes.

* `proxybuffer_high` - (Optional,type `int`) Specifies the proxy buffer level, in bytes, at which the receive window is closed.

* `congestion_control` - (Optional,type `string`) Specifies the algorithm to use to share network resources among competing users to reduce congestion. The default is High Speed.

* `initial_congestion_windowsize` - (Optional,type `int`) Specifies the initial congestion window size for connections to this destination. Actual window size is this value multiplied by the MSS (Maximum Segment Size) for the same connection. The default is 10. Valid values range from 0 to 64.

* `delayed_acks` - (Optional,type `string`) Specifies, when checked (enabled), that the system can send fewer than one ACK (acknowledgment) segment per data segment received. By default, this setting is enabled.

* `nagle` - (Optional,type `string`) Specifies whether the system applies Nagle's algorithm to reduce the number of short segments on the network.If you select Auto, the system determines whether to use Nagle's algorithm based on network conditions. By default, this setting is disabled.

* `early_retransmit` - (Optional,type `string`) Enabling this setting allows TCP to assume a packet is lost after fewer than the standard number of duplicate ACKs, if there is no way to send new data and generate more duplicate ACKs.

* `tailloss_probe` - (Optional,type `string`) Enabling this setting allows TCP to send a probe segment to trigger fast recovery instead of recovering a loss via a retransmission timeout,By default, this setting is enabled.

* `timewait_recycle` - (Optional,type `string`) Using this setting enabled, the system can recycle a wait-state connection immediately upon receipt of a new connection request instead of having to wait until the connection times out of the wait state. By default, this setting is enabled.

* `fast_open` - (Optional,type `string`) When enabled, permits TCP Fast Open, allowing properly equipped TCP clients to send data with the SYN packet. Default is `enabled`. If `fast_open` set to `enabled`, argument `verified_accept` can't be set to `enabled`.

* `verified_accept` - (Optional,type `string`) Specifies, when checked (enabled), that the system can actually communicate with the server before establishing a client connection. To determine this, the system sends the server a SYN packet before responding to the client's SYN with a SYN-ACK. When unchecked, the system accepts the client connection before selecting a server to talk to. By default, this setting is `disabled`.

* `deferred_accept` - (Optional,type `string`) Specifies, when enabled, that the system defers allocation of the connection chain context until the client response is received. This option is useful for dealing with 3-way handshake DOS attacks. The default value is disabled.

## Importing
An existing tcp profile can be imported into this resource by supplying tcp profile Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_profile_tcp.tcp-lan-profile-import /Common/test-tcp-lan-profile
```
