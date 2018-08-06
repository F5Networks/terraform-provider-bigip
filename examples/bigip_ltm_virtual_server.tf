provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_virtual_server" "http" {
	pool = "/Common/App1_pool"
        name = "/Common/App1_vs_http"
	destination = "100.1.1.100"
	port = 80
	source_address_translation = "automap"
}

resource "bigip_ltm_pool"  "App1_pool" {
         name = "/Common/App1_pool"
        load_balancing_mode = "round-robin"
        monitors = ["/Common/App1_monitor"]
        allow_snat = true
       }
 
resource "bigip_ltm_pool_attachment" "attach_node1" {
        pool = "/Common/App1_pool"
            node = "/Common/11.1.1.101:80"
            depends_on = ["bigip_ltm_pool.App1_pool"]
 
}
resource "bigip_ltm_pool_attachment" "attach_node2" {
        pool = "/Common/App1_pool"
            node = "/Common/11.1.1.102:80"
            depends_on = ["bigip_ltm_pool.App1_pool"]
 
}
