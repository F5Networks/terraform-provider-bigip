/* Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0. 
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

provider "bigip" {
  address = "x.x.x.x"
  username = "xxxx"
  password = "xxxxx"
}


// tenant_name is used to set the identity of as3 resource which is unique for resource.
resource "bigip_as3"  "as3-example1" {
     as3_json = "${file("example1.json")}" 
     tenant_name = "as3"
 }

