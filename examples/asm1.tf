provider "bigip" {
  address = "10.192.74.61"
  username = "admin"
  password = "admin"
}


resource "bigip_ltm_provision" "provision-ilx" {
  name = "/Common/ilx"
  fullPath  = "ilx"
  cpuRatio = 0
  diskRatio = 0
  level = "nominal"
  memoryRatio = 0
}

