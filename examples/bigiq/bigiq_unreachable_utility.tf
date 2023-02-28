provider "bigip" {
  address  = "XXX.XXX.XXX.XXXX"
  username = "XXXXXX"
  password = "XXXXXXX"
}
resource "bigip_common_license_manage_bigiq" "test_example" {
  bigiq_address    = "xxx.xxx.xxx.xxx"
  bigiq_user       = "xxxx"
  bigiq_password   = "xxxxx"
  license_poolname = "utility_pool_name"
  unit_of_measure  = "yearly"
  assignment_type  = "UNREACHABLE"
  mac_address      = "FA:16:3E:1B:6D:32"
  hypervisor       = "azure"
}