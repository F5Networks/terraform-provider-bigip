provider "bigip" {
  address  = "XXX.XXX.XXX.XXXX"
  username = "XXXXXX"
  password = "XXXXXXX"
}
resource "bigip_common_license_manage_bigiq" "test_example" {
  bigiq_address    = "xxx.xxx.xxx.xxx"
  bigiq_user       = "xxxx"
  bigiq_password   = "xxxxx"
  license_poolname = "purchased_pool_name"
  assignment_type  = "MANAGED"
}
