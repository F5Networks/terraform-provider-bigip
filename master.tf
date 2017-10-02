provider "bigip" {
  address = "10.0.1.79"
  username = "admin"
  password = "admin"
}
resource "bigip_syslog" "sys1" {
  auth_privfrom = "notice"
  remote_servers = {
  	name = "/Common/test"
  	host = "myserver1.com"
    remote_port = 514
 }
}
