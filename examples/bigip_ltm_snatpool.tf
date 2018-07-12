provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_ltm_snatpool" "snatpool_sanjose" {
  name = "/Common/snatpool_sanjose"
  members = ["191.1.1.1","194.2.2.2"]
}

