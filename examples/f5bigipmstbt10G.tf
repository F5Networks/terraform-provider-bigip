provider "bigip" {
/* BIG-IQ CM used to license BIG-IP devices */ 
  address = "10.192.74.92"
  username = "admin"
  password = "admin"
}

resource "bigip_license_f5bigmstbt10G" "San_Jose_f5bigipmstbt10G" {
/* BIG-IP at San Jose Licensed */  
  device_address = "10.192.74.73"
  username = "admin"
  password = "admin"
  unitof_measure = "yearly"
}
