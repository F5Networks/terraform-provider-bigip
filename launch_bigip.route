provider "bigip" {
  address = "10.192.74.61"
  username = "admin"
  password = "admin"
}

resource "bigip_route" "route2" {
  name = "sanjay-route2"
  network = "10.10.10.0/24"
  gw      = "1.1.1.2"
}

