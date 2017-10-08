provider "bigip" {
  address = "10.0.1.79"
  username = "admin"
  password = "admin"
}

resource "bigip_datagroup" "datagroup1" {
  name = "sanjosedatagr"
  type = "string"
  records  {
   name = "abc.com"
   data = "pool1"
   }
}
