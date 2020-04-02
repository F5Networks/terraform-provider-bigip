---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_monitor"
sidebar_current: "docs-bigip-resource-monitor-x"
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

* `name` (Required) Name of the monitor

* `parent` - (Required) Existing LTM monitor to inherit from

* `interval` - (Optional) Check interval in seconds

* `timeout` - (Optional) Timeout in seconds

* `send` - (Optional) Request string to send

* `receive` - (Optional) Expected response string

* `receive_disable` - (Optional)

* `reverse`  - (Optional)

* `transparent` - (Optional)

* `manual_resume` - (Optional)

* `ip_dscp` - (Optional)

* `time_until_up` - (Optional)

* `database` - (Optional) Specifies the database in which the user is created

* `destination` - (Optional) Specify an alias address for monitoring

* `username` - (Optional) Specifies the user name if the monitored target requires authentication

* `password` - (Optional) Specifies the password if the monitored target requires authentication 

* `compatibility` -  (Optional) Specifies, when enabled, that the SSL options setting (in OpenSSL) is set to ALL. Accepts 'enabled' or 'disabled' values, the default value is 'enabled'.

* `filename` - (Optional) Specifies the full path and file name of the file that the system attempts to download. The health check is successful if the system can download the file.

* `mode` - (Optional) Specifies the data transfer process (DTP) mode. The default value is passive. The options are passive (Specifies that the monitor sends a data transfer request to the FTP server. When the FTP server receives the request, the FTP server then initiates and establishes the data connection.) and active (Specifies that the monitor initiates and establishes the data connection with the FTP server.).
