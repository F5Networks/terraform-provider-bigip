provider "bigip" {
  address = "10.192.74.80"
  username = "admin"
  password = "admin"
}

resource "bigip_license" "lic13" {
  deviceAddress = "10.192.74.55"
  username = "admin"
  password = "admin"
}
