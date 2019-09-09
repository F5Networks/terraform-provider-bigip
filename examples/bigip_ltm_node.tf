provider "bigip" {
  address = "10.192.74.68"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_node" "node" {
  name = "/Common/terraform_node1"
  address = "10.10.10.10"
}

