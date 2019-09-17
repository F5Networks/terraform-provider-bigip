/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_ltm_snat" "snat_list" {
 name = "NewSnatList"
 translation = "136.1.1.1"
 origins =  {name = "2.2.2.2"}
 origins =  {name = "3.3.3.3"} 
}

