provider "bigip" {
address = "10.192.74.73"
username = "admin"
password = "admin"

}
module  "sjvlan1" {
  source = "./vlanmodule"
  name = "/Common/intvlan"
  tag = 101
  vlanport = "1.1"
  tagged = true
 }

resource "bigip_ltm_selfip" "selfip" {

        name = "/Common/InternalselfIP"
        ip = "100.1.1.1/24"
        vlan = "/Common/intvlan"
        depends_on = ["module.sjvlan1"]
        }



