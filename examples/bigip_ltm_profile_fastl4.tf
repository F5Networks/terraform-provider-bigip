/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxx"
}

resource "bigip_ltm_profile_fastl4" "sjfastl4profile" {
  name                   = "/Common/sjfastl4profile"
  partition              = "Common"
  defaults_from          = "/Common/fastL4"
  client_timeout         = 40
  explicitflow_migration = "enabled"
  hardware_syncookie     = "enabled"
  idle_timeout           = "200"
  iptos_toclient         = "pass-through"
  iptos_toserver         = "pass-through"
  keepalive_interval     = "disabled" //This cannot take enabled
}

