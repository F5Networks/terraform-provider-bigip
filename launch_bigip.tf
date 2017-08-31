provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_datagroup" "datagroup1" {
  name = "dgx8"
  type = "string"
  records  {
   name = "xyx.f5.com"
   data = "pool100"
   }
}
