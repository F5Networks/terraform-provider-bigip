provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}

resource "bigip_datagroup" "datagroup1" {
  name = "dgx2"
  type = "string"
  records  {
   name = "abc.com"
   data = "pool1"
   }
}
