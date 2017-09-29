provider "bigip" {
  address = "10.0.1.79"
  username = "admin"
  password = "admin"
}
resource "bigip_ltm_provision" "provision-afm" {
  name = "/Common/afm"
  full_path  = "afm"
  cpu_ratio = 0
  disk_ratio = 0
  level = "nominal"
  memory_ratio = 0
}

resource "bigip_ltm_ntp" "ntp1" {
  	description = "/Common/NTP1"
  	servers = ["time.google.com"]
  	timezone = "America/Los_Angeles"
	depends_on = ["bigip_ltm_provision.provision-afm"]
}

resource "bigip_ltm_dns" "dns1" {
   	description = "/Common/DNS1"
   	name_servers = ["8.8.8.8"]
   	numberof_dots = 2
   	search = ["f5.com"]
   	depends_on = ["bigip_ltm_provision.provision-afm"]
}
resource "bigip_ltm_vlan" "vlan1" {
	name = "/Common/internal"
	tag = 101
	interfaces = {
		vlanport = 1.2,
		tagged = false
	}	

        depends_on = ["bigip_ltm_provision.provision-afm"]
}

resource "bigip_ltm_vlan" "vlan2" {
        name = "/Common/external"
        tag = 102
        interfaces = {
                vlanport = 1.1,
                tagged = false
        }

        depends_on = ["bigip_ltm_provision.provision-afm"]
}

resource "bigip_ltm_selfip" "selfip1" {
	name = "/Common/internalselfIP"
	ip = "11.1.1.1/24"
	vlan = "/Common/internal"
	depends_on = ["bigip_ltm_vlan.vlan1"]
        depends_on = ["bigip_ltm_provision.provision-afm"]
	}

resource "bigip_ltm_selfip" "selfip2" {
        name = "/Common/externalselfIP"
        ip = "100.1.1.1/24"
        vlan = "/Common/external"
        depends_on = ["bigip_ltm_vlan.vlan2"]
        depends_on = ["bigip_ltm_provision.provision-afm"]
        }


resource "bigip_ltm_monitor" "monitor" {
        name = "/Common/terraform_monitor"
        parent = "/Common/http"
        send = "GET /some/path\r\n"
        timeout = "999"
        interval = "999"
        depends_on = ["bigip_ltm_provision.provision-afm"]
}

resource "bigip_ltm_pool"  "pool" {
        name = "/Common/terraform-pool"
        load_balancing_mode = "round-robin"
        nodes = ["11.1.1.101:80", "11.1.1.102:80"]
        monitors = ["/Common/terraform_monitor"]
        allow_snat = true
        depends_on = ["bigip_ltm_provision.provision-afm"]
}

resource "bigip_ltm_virtual_server" "http" {
	pool = "/Common/terraform-pool"
        name = "/Common/terraform_vs_http"
	destination = "100.1.1.100"
	port = 80
	source_address_translation = "automap"
        depends_on = ["bigip_ltm_pool.pool"]
}



