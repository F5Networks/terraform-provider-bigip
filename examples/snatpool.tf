provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_ltm_snatpool" "snatpool_sanjose" {
  // Modification APIs are not supported you cannot do PUT or PAtch for this resource only Delete and Create
  name = "/Common/snatpool_sanjose"
  members = ["191.1.1.1","194.2.2.2"]
}

