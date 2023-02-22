provider "bigip" {
  address  = "XXX.XXX.XXX.XXXX"
  username = "XXXXXX"
  password = "XXXXXXX"
}
resource "bigip_common_license_manage_bigiq" "test_example" {
  bigiq_address    = "xxx.xxx.xxx.xxx"
  bigiq_user       = "xxxx"
  bigiq_password   = "xxxxx"
  license_poolname = "regkey_pool_name"
  assignment_type  = "UNMANAGED"
  key              = "W8368-38939-99443-37082-0654410"
}