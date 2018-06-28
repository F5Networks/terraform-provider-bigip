provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_pool_attachment" "attach_node" {
        pool = "/Common/terraform-pool"
	node = "/Common/11.1.1.101:80"
	depends_on = ["bigip_ltm_pool.pool"]

}


