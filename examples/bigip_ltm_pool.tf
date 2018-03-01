provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_pool"  "pool" {
        name = "/Common/terraform-pool"
        load_balancing_mode = "round-robin"
        nodes = ["11.1.1.101:80", "11.1.1.102:80"]
        allow_snat = "yes"
        allow_nat = "yes"
}

