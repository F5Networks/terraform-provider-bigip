/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

provider "bigip" {
  address  = "x.x.x.x"
  username = "xxxx"
  password = "xxxx"
}

// Using  provisioner to download and install do rpm on bigip, pass arguments as BIG-IP IP address, credentials 
// Use this provisioner for first time to download and install do rpm on bigip
resource "null_resource" "install_do" {
  provisioner "local-exec" {
    command = "./install-do-rpm.sh x.x.x.x xxxx:xxxx"
  }
}

// config_name is used to set the identity of do resource which is unique for resource.
resource "bigip_do" "do-example1" {
  do_json     = file("example1.json")
  config_name = "sample_test"
  depends_on  = ["null_resource.install_do"]

}

