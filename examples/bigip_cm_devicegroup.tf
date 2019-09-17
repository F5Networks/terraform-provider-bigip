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


resource "bigip_cm_devicegroup" "my_new_devicegroup"

        {
            name = "deadlygroup"
            auto_sync = "enabled"
            full_load_on_sync = "true"
            type = "sync-only"
            device  { name = "bigip1.cisco.com"}
            device  { name = "bigip200.f5.com"}
        }
