/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

provider "bigip" {
  address = "x.x.x.x"
  username = "xxxx"
  password = "xxxxx"
}


// Using  provisioner to download and install do rpm on bigip, pass arguments as BIG-IP IP address, credentials
// Use this provisioner for first time to download and install do rpm on bigip
resource "null_resource" "install_as3" {
  provisioner "local-exec" {
    command = "./install-as3-rpm.sh x.x.x.x xxxx:xxxx"
  }
}

// config_name is used to set the identity of as3 resource which is unique for resource.
resource "bigip_as3"  "as3-example1" {
     as3_json = "${file("example1.json")}" 
     config_name = "sample_test"
     depends_on = ["null_resource.install_as3"]
 }


