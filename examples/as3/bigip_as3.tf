provider "bigip" {
  address = "x.x.x.x"
  username = "xxxx"
  password = "xxxxx"
}


// tenant_name is used to set the identity of as3 resource which is unique for resource.
resource "bigip_as3"  "as3-example1" {
     as3_json = "${file("example1.json")}" 
     tenant_name = "as3"
 }

