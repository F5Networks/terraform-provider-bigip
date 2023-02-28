/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
//We can use null_resource to  deploy As3 templates, below is simple example to install the as3 rpm and another resource which deploys the example1.json ( which has the http VS configuration) More details on As3 please refer to https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/


provider "bigip" {
  address  = "X.X.X.X"
  username = "xxxx"
  password = "xxxx"
}


// Using  provisioner to install as3 rpm on bigip pass arguments as BIG-IP IP address, credentials and name of the rpm 
resource "null_resource" "install_as3" {
  provisioner "local-exec" {
    command = "sh install_as3.sh X.X.X.X  admin:pass f5-appsvcs-3.9.0-3.noarch.rpm"
  }
}
// The below null resource can be used to deploy HTTP script as3_http.sh which uses JSON payload example1.josn 
resource "null_resource" "deploy_as3_http" {
  provisioner "local-exec" {
    command = "sh as3_http.sh"
  }
  depends_on = ["null_resource.install_as3"]

}

