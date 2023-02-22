provider "bigip" {
  address  = "xx.xx.xx.xxx"
  username = "xxxx"
  password = "xxxxxx"
}
resource "bigip_as3" "as3-example1" {
  as3_json = file("as3_example1.json")
}

resource "bigip_as3" "as3-example2" {
  as3_json = file("as3_example2.json")
}

