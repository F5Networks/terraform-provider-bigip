---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_http2"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_http2 resource
---

# bigip\_ltm\_profile_http2

`bigip_ltm_profile_http2` Configures a custom profile_http2 for use by health checks.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage

```hcl

resource "bigip_ltm_profile_http2" "nyhttp2" {
  name                              = "/Common/test-profile-http2"
  frame_size                        = 2021
  receive_window                    = 31
  write_size                        = 16380
  header_table_size                 = 4092
  include_content_length            = "enabled"
  enforce_tls_requirements          = "enabled"
  insert_header                     = "disabled"
  concurrent_streams_per_connection = 30
  connection_idle_timeout           = 100
  activation_modes                  = ["always"]
}

#Child Profile which inherits parent http2 profile

resource "bigip_ltm_profile_http2" "nyhttp2-child" {
  name          = "/Common/test-profile-http2-child"
  defaults_from = bigip_ltm_profile_http2.nyhttp2.name
}

```      

## Argument Reference

* `name` (Required,`type string`) Name of Profile should be full path.The full path is the combination of the `partition + profile name`,For example `/Common/test-http2-profile`.

* `defaults_from` - (Optional,`type string`) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `concurrent_streams_per_connection` - (Optional,`type int`) Specifies how many concurrent requests are allowed to be outstanding on a single HTTP/2 connection.

* `connection_idle_timeout` - (Optional,`type int`) Specifies the number of seconds that a connection is idle before the connection is eligible for deletion.

* `insert_header` - (Optional,`type string`) This setting specifies whether the BIG-IP system should add an HTTP header to the HTTP request to show that the request was received over HTTP/2, Allowed Values : `"enabled"/"disabled"` [ Default: `"disabled"`].

* `insert_header_name` - (Optional,`type string`) This setting specifies the name of the header that the BIG-IP system will add to the HTTP request when the Insert Header is enabled.

* `header_table_size` - (Optional) The size of the header table, in KB, for the HTTP headers that the HTTP/2 protocol compresses to save bandwidth.

* `enforce_tls_requirements` - (Optional,`type string`) Enable or disable enforcement of TLS requirements,Allowed Values : `"enabled"/"disabled"` [Default:`"enabled"`].

* `frame_size` - (Optional,`type int`) The size of the data frames, in bytes, that the HTTP/2 protocol sends to the client. `Default: 2048`.

* `receive_window` - (Optional,`type int`) The flow-control size for upload streams, in KB. `Default: 32`.

* `write_size` - (Optional,`type int`) The total size of combined data frames, in bytes, that the HTTP/2 protocol sends in a single write function. `Default: 16384`".

* `activation_modes` - (Optional) This setting specifies the condition that will cause the BIG-IP system to handle an incoming connection as an HTTP/2 connection, Allowed values : `[“alpn”]` (or) `[“always”]`.
