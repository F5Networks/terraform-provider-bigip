---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_monitor"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_monitor resource
---

# bigip\_ltm\_monitor

`bigip_ltm_monitor` Configures a custom monitor for use by health checks.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage


```hcl

resource "bigip_ltm_monitor" "monitor" {
  name        = "/Common/terraform_monitor"
  parent      = "/Common/http"
  send        = "GET /some/path\r\n"
  timeout     = "999"
  interval    = "999"
  destination = "1.2.3.4:1234"
}

resource "bigip_ltm_monitor" "test-https-monitor" {
  name        = "/Common/terraform_monitor"
  parent      = "/Common/http"
  ssl_profile = "/Common/serverssl"
  send        = "GET /some/path\r\n"
  timeout     = "999"
  interval    = "999"
}

resource "bigip_ltm_monitor" "test-ftp-monitor" {
  name          = "/Common/ftp-test"
  parent        = "/Common/ftp"
  interval      = 5
  time_until_up = 0
  timeout       = 16
  destination   = "*:8008"
  filename      = "somefile"
}

resource "bigip_ltm_monitor" "test-postgresql-monitor" {
  name     = "/Common/test-postgresql-monitor"
  parent   = "/Common/postgresql"
  send     = "SELECT 'Test';"
  receive  = "Test"
  interval = 5
  timeout  = 16
  username = "abcd"
  password = "abcd1234"
}
```      

## Argument Reference

* `name` ((Required,type `string`) Specifies the Name of the LTM Monitor.Name of Monitor should be full path,full path is the combination of the `partition + monitor name`,For ex:`/Common/test-ltm-monitor`.

* `parent` - (Required,type `string`)  Parent monitor for the system to use for setting initial values for the new monitor.

* `interval` - (Optional,type `int`) Specifies, in seconds, the frequency at which the system issues the monitor check when either the resource is down or the status of the resource is unknown. The default is `5`

* `up_interval` - (Optional,type `int`) Specifies the interval for the system to use to perform the health check when a resource is up. The default is `0(Disabled)`

* `timeout` - (Optional,type `int`) Specifies the number of seconds the target has in which to respond to the monitor request. The default is `16` seconds

* `send` - (Optional,type `string`) Specifies the text string that the monitor sends to the target object.

* `receive` - (Optional,type `string`) Specifies the regular expression representing the text string that the monitor looks for in the returned resource.

* `receive_disable` - (Optional,type `string`) The system marks the node or pool member disabled when its response matches Receive Disable String but not Receive String.

* `reverse`  - (Optional,type `string`) Instructs the system to mark the target resource down when the test is successful.

* `transparent` - (Optional,type `string`) Specifies whether the monitor operates in transparent mode.

* `manual_resume` - (Optional,type `string`) Specifies whether the system automatically changes the status of a resource to Enabled at the next successful monitor check.

* `ip_dscp` - (Optional,type `int`) Displays the differentiated services code point (DSCP).The default is `0 (zero)`.

* `time_until_up` - (Optional,type `int`) Specifies the number of seconds to wait after a resource first responds correctly to the monitor before setting the resource to up.

* `database` - (Optional) Specifies the database in which the user is created

* `destination` - (Optional,type `string`) Specify an alias address for monitoring

* `adaptive` - (Optional,type `string`) Specifies whether adaptive response time monitoring is enabled for this monitor. The default is `disabled`.

* `adaptive_limit` - (Optional,type `int`) Specifies the absolute number of milliseconds that may not be exceeded by a monitor probe, regardless of Allowed Divergence.

* `username` - (Optional,type `string`) Specifies the user name if the monitored target requires authentication

* `password` - (Optional,type `string`) Specifies the password if the monitored target requires authentication 

* `compatibility` -  (Optional,type `string`) Specifies, when enabled, that the SSL options setting (in OpenSSL) is set to ALL. Accepts 'enabled' or 'disabled' values, the default value is 'enabled'.

* `filename` - (Optional,type `string`) Specifies the full path and file name of the file that the system attempts to download. The health check is successful if the system can download the file.

* `mode` - (Optional,type `string`) Specifies the data transfer process (DTP) mode. The default value is passive. The options are passive (Specifies that the monitor sends a data transfer request to the FTP server. When the FTP server receives the request, the FTP server then initiates and establishes the data connection.) and active (Specifies that the monitor initiates and establishes the data connection with the FTP server.).

* `ssl_profile` - (Optional,type `string`) Specifies the ssl profile for the monitor. It only makes sense when the parent is `/Common/https`
