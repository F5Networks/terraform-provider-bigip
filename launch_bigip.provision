provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_ltm_provision" "provision1" {
  name = "/Common/afm"
  full_path  = "afm"
  cpu_ratio = 0
  disk_ratio = 0
  level = "nominal"
  memory_ratio = 0
}

