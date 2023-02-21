provider "bigip" {
  address  = "xx.xxx.xx.xxx"
  username = "xxxxx"
  password = "xxxxxxxx"
}

resource "bigip_bigiq_as3" "exampletask" {
  bigiq_address  = "xx.xx.xxx.xx"
  bigiq_user     = "xxxxx"
  bigiq_password = "xxxxxxxxx"
  as3_json       = file("bigiq_example.json")
}
