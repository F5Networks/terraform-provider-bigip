provider "bigip" {
/* BIG-IQ CM used to license BIG-IP devices */ 
  address = "10.192.74.80"
  username = "admin"
  password = "admin"
}

resource "bigip_license" "San_Jose" {
/* BIG-IP at San Jose Licensed */  
  device_address = "10.192.74.55"
  username = "admin"
  password = "admin"
}
resource "bigip_license" "New_York" {
/* BIG-IP at New York  Licensed */  
  device_address = "10.192.74.61"
  username = "admin"
  password = "admin"
}
resource "bigip_license" "San_Antonio" {
/* BIG-IP at San Antonio Licensed */  
  device_address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_license" "Seattle" {
/* BIG-IP at Seattle Licensed */
  device_address = "10.192.74.81"
  username = "admin"
  password = "admin"
}

