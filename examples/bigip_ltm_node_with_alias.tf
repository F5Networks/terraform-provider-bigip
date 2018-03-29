provider "bigip" {
  address = "10.192.74.73"
  alias = "east"
  username = "admin"
  password = "admin"
}

provider "bigip" {
   alias = "west"
   address = "10.192.74.68"
   username = "admin"
   password = "admin"
}


resource "bigip_ltm_node" "node_west" {
  name = "/Common/terraform_node1"
  provider = "bigip.west"
  address = "1.1.1.1"
  state = "user-up"
}

resource "bigip_ltm_node" "node_east" {
  name = "/Common/terraform_node1"
  provider = "bigip.east"
  address = "1.1.1.1"
  state = "user-down"
}



