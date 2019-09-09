provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_monitor" "monitor" {
        name = "/Common/terraform_monitor"
        parent = "/Common/http"
        send = "GET /some/path\r\n"
        timeout = "999"
        interval = "999"
        depends_on = ["bigip_sys_provision.provision-afm"]
}

