provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}
resource "bigip_syslog" "sys1" {
  remote_servers = {
  	name = "/Common/sanjay"
  	host = "myserver1.com"
    remote_port = 514
 }
}

