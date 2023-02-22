provider "bigip" {
  address  = "XXX.XXX.XXX.XXXX"
  username = "XXXXXX"
  password = "XXXXXXX"
}

resource "bigip_common_license_manage_bigiq" "test_example" {
  bigiq_address    = "xxx.xxx.xxx.xxx"
  bigiq_user       = "xxxx"
  bigiq_password   = "xxxx"
  license_poolname = "utility_license_name"
  assignment_type  = "UNMANAGED"
  unit_of_measure  = "yearly"
}